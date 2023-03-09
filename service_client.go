package weshnet

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	"berty.tech/weshnet/pkg/grpcutil"
	"berty.tech/weshnet/pkg/protocoltypes"
)

const ClientBufferSize = 4 * 1024 * 1024

type ServiceClient interface {
	protocoltypes.ProtocolServiceClient

	io.Closer
}

type serviceClient struct {
	ServiceClient // inehrit from client

	service Service
	server  *grpc.Server
}

func (c *serviceClient) Close() error {
	c.server.GracefulStop()     // gracefuly stop grpc server
	_ = c.ServiceClient.Close() // close client and discard error
	return c.service.Close()    // return real service error
}

func NewServiceClient(opts Opts) (ServiceClient, error) {
	svc, err := NewService(opts)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()
	c, err := NewClientFromService(context.Background(), s, svc)
	if err != nil {
		return nil, fmt.Errorf("uanble to create client from server: %w", err)
	}

	return &serviceClient{
		ServiceClient: c,
		server:        s,
		service:       svc,
	}, nil
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
