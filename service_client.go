package weshnet

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"berty.tech/weshnet/v2/pkg/grpcutil"
	"berty.tech/weshnet/v2/pkg/ipfsutil"
	ipfs_mobile "berty.tech/weshnet/v2/pkg/ipfsutil/mobile"
	"berty.tech/weshnet/v2/pkg/logutil"
	"berty.tech/weshnet/v2/pkg/protocoltypes"
)

const (
	defaultLoggingFiltersKey   = ":default:"
	defaultLoggingFiltersValue = "info+:bty.* error+:*,-ipfs*,-*.tyber"
)

type ServiceClient interface {
	protocoltypes.ProtocolServiceClient

	io.Closer
}

// NewServiceClient initializes a new ServiceClient using the opts.
// If opts.RootDatastore is nil and opts.DatastoreDir is "" or InMemoryDirectory, then set
// opts.RootDatastore to an in-memory data store. Otherwise, if opts.RootDatastore is nil then set
// opts.RootDatastore to a persistent data store at opts.DatastoreDir .
func NewServiceClient(opts Opts) (ServiceClient, error) {
	var err error

	var cleanupLogger func()
	if opts.Logger == nil {
		if opts.Logger, cleanupLogger, err = setupDefaultLogger(); err != nil {
			return nil, fmt.Errorf("unable to setup logger: %w", err)
		}
	}

	svc, err := NewService(opts)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	c, err := NewClientFromService(ctx, s, svc)
	if err != nil {
		return nil, fmt.Errorf("uanble to create client from server: %w", err)
	}

	return &serviceClient{
		ServiceClient: c,
		server:        s,
		service:       svc,
		cleanup:       cleanupLogger,
	}, nil
}

// NewInMemoryServiceClient creates a new in-memory Wesh protocol service and returns a gRPC
// ServiceClient which uses a direct in-memory connection. When finished, you must call Close().
// This creates a new Wesh account where the key store is in memory. (If you don't
// export the data then it is lost when you call Close(). ) The IPFS node, cached data,
// and configuration are also in memory.
func NewInMemoryServiceClient() (ServiceClient, error) {
	var opts Opts
	opts.DatastoreDir = InMemoryDirectory
	return NewServiceClient(opts)
}

// NewPersistentServiceClient creates a Wesh protocol service using persistent storage files in the
// directory given by the directory path. If the directory doesn't exist, this creates it with files
// of a new Wesh account and peer identity. (If the directory doesn't exist, this will create it only
// if the parent directory exists. Otherwise you must first create the parent directories.) However,
// if the persistent storage files already exist, then this opens them to use the existing Wesh
// account and peer identity. This returns a gRPC ServiceClient which uses a direct in-memory
// connection. When finished, you must call Close().
func NewPersistentServiceClient(path string) (ServiceClient, error) {
	var opts Opts

	opts.DatastoreDir = path

	repo, err := ipfsutil.LoadRepoFromPath(path)
	if err != nil {
		return nil, err
	}

	mrepo := ipfs_mobile.NewRepoMobile(path, repo)
	mnode, err := ipfsutil.NewIPFSMobile(context.TODO(), mrepo, &ipfsutil.MobileOptions{})
	if err != nil {
		return nil, err
	}

	opts.IpfsCoreAPI, err = ipfsutil.NewExtendedCoreAPIFromNode(mnode.IpfsNode)
	if err != nil {
		return nil, err
	}

	var cleanupLogger func()
	if opts.Logger, cleanupLogger, err = setupDefaultLogger(); err != nil {
		return nil, fmt.Errorf("uanble to setup logger: %w", err)
	}

	cl, err := NewServiceClient(opts)
	if err != nil {
		return nil, err
	}

	return &persistentServiceClient{
		ServiceClient: cl,
		cleanup:       cleanupLogger,
	}, nil
}

const ClientBufferSize = 4 * 1024 * 1024

type serviceClient struct {
	ServiceClient // inehrit from client

	service Service
	server  *grpc.Server
	cleanup func()
}

type persistentServiceClient struct {
	ServiceClient
	cleanup func()
}

func (p *persistentServiceClient) Close() error {
	err := p.ServiceClient.Close()

	if p.cleanup != nil {
		p.cleanup()
	}

	return err
}

func (c *serviceClient) Close() (err error) {
	c.server.GracefulStop()     // gracefully stop grpc server
	_ = c.ServiceClient.Close() // close client and discard error

	err = c.service.Close()

	if c.cleanup != nil {
		c.cleanup()
	}

	return // return real service error
}

type client struct {
	protocoltypes.ProtocolServiceClient

	l  *grpcutil.BufListener
	cc *grpc.ClientConn
}

func (c *client) Close() error {
	err := c.cc.Close()
	_ = c.l.Close()
	return err
}

func NewClientFromService(ctx context.Context, s *grpc.Server, svc Service, opts ...grpc.DialOption) (ServiceClient, error) {
	bl := grpcutil.NewBufListener(ClientBufferSize)
	cc, err := bl.NewClientConn(ctx, opts...)
	if err != nil {
		return nil, err
	}

	protocoltypes.RegisterProtocolServiceServer(s, svc)
	go func() {
		// we dont need to log the error
		_ = s.Serve(bl)
	}()

	return &client{
		ProtocolServiceClient: protocoltypes.NewProtocolServiceClient(cc),
		cc:                    cc,
		l:                     bl,
	}, nil
}

func setupDefaultLogger() (logger *zap.Logger, cleanup func(), err error) {
	// setup log from env
	if logfilter := os.Getenv("WESHNET_LOG_FILTER"); logfilter != "" {
		if logfilter == defaultLoggingFiltersKey {
			logfilter = defaultLoggingFiltersValue
		}

		s := logutil.NewStdStream(logfilter, "color", os.Stderr.Name())
		return logutil.NewLogger(s)
	}

	return zap.NewNop(), func() {}, nil
}
