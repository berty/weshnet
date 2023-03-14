package weshnet

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/dgraph-io/badger/v2/options"
	"github.com/ipfs/go-datastore"
	badger "github.com/ipfs/go-ds-badger2"
	"google.golang.org/grpc"

	"berty.tech/weshnet/pkg/grpcutil"
	"berty.tech/weshnet/pkg/ipfsutil"
	ipfs_mobile "berty.tech/weshnet/pkg/ipfsutil/mobile"
	"berty.tech/weshnet/pkg/protocoltypes"
)

type ServiceClient interface {
	protocoltypes.ProtocolServiceClient

	io.Closer
}

func NewServiceClient(opts Opts) (ServiceClient, error) {
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
	}, nil
}

func NewInMemoryServiceClient() (ServiceClient, error) {
	var opts Opts
	opts.DatastoreDir = InMemoryDirectory
	return NewServiceClient(opts)
}

func NewPersistantServiceClient(path string) (ServiceClient, error) {
	var opts Opts

	bopts := badger.DefaultOptions
	bopts.ValueLogLoadingMode = options.FileIO

	ds, err := badger.NewDatastore(path, &bopts)
	if err != nil {
		return nil, fmt.Errorf("unable to init badger datastore: %w", err)
	}

	repo, err := ipfsutil.LoadRepoFromPath(path)
	if err != nil {
		return nil, err
	}

	mrepo := ipfs_mobile.NewRepoMobile(path, repo)
	mnode, err := ipfsutil.NewIPFSMobile(context.TODO(), mrepo, &ipfsutil.MobileOptions{
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	})
	if err != nil {
		return nil, err
	}

	opts.IpfsCoreAPI, err = ipfsutil.NewExtendedCoreAPIFromNode(mnode.IpfsNode)
	if err != nil {
		return nil, err
	}

	opts.RootDatastore = ds

	cl, err := NewServiceClient(opts)
	if err != nil {
		return nil, err
	}

	return &persistantServiceClient{
		ServiceClient: cl,
		ds:            ds,
	}, nil
}

const ClientBufferSize = 4 * 1024 * 1024

type serviceClient struct {
	ServiceClient // inehrit from client

	service Service
	server  *grpc.Server
}

type persistantServiceClient struct {
	ServiceClient
	ds datastore.Batching
}

func (p *persistantServiceClient) Close() error {
	err := p.ServiceClient.Close()

	if dserr := p.ds.Close(); err == nil && dserr != nil {
		// only return ds error if no error have been catch earlier
		err = fmt.Errorf("unable to close datastore: %w", dserr)
	}
	return err
}

func (c *serviceClient) Close() error {
	c.server.GracefulStop()     // gracefully stop grpc server
	_ = c.ServiceClient.Close() // close client and discard error
	return c.service.Close()    // return real service error
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
		switch err := s.Serve(bl); err {
		case nil, grpc.ErrServerStopped: // ok
		default:
			panic(err)
		}
	}()

	return &client{
		ProtocolServiceClient: protocoltypes.NewProtocolServiceClient(cc),
		cc:                    cc,
		l:                     bl,
	}, nil
}
