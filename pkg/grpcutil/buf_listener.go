package grpcutil

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type BufListener struct {
	*bufconn.Listener
}

func NewBufListener(sz int) *BufListener {
	return &BufListener{
		Listener: bufconn.Listen(sz),
	}
}

func (bl *BufListener) dialer(context.Context, string) (net.Conn, error) {
	return bl.Dial()
}

func (bl *BufListener) NewClientConn(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	mendatoryOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(bl.dialer), // set pipe dialer
	}

	return grpc.DialContext(ctx, "buf", append(opts, mendatoryOpts...)...)
}
