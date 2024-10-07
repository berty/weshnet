package weshnet

import (
	"context"
	"fmt"
	"io"
	"time"

	ds "github.com/ipfs/go-datastore"
	ds_sync "github.com/ipfs/go-datastore/sync"
	coreiface "github.com/ipfs/kubo/core/coreiface"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"berty.tech/weshnet/v2/pkg/errcode"
	"berty.tech/weshnet/v2/pkg/grpcutil"
	"berty.tech/weshnet/v2/pkg/outofstoremessagetypes"
	"berty.tech/weshnet/v2/pkg/protocoltypes"
	"berty.tech/weshnet/v2/pkg/secretstore"
)

type OOSMService interface {
	outofstoremessagetypes.OutOfStoreMessageServiceServer
}

var _ OOSMService = (*oosmService)(nil)

type oosmService struct {
	logger        *zap.Logger
	rootDatastore ds.Datastore
	secretStore   secretstore.SecretStore

	outofstoremessagetypes.UnimplementedOutOfStoreMessageServiceServer
}

type OOSMServiceClient interface {
	outofstoremessagetypes.OutOfStoreMessageServiceClient

	io.Closer
}

type oosmServiceClient struct {
	OOSMServiceClient

	service OOSMService
	server  *grpc.Server
}

type OOSMOption func(*oosmService) error

// NewOutOfStoreMessageServiceClient creates a new Wesh protocol service and returns a gRPC
// ServiceClient which uses a direct in-memory connection. When finished, you must call Close().
// This opens or creates a Wesh account where the datastore location is specified by the path argument.
// The service will not start any network stuff, it will only use the filesystem to store or get data.
func NewOutOfStoreMessageServiceClient(opts ...OOSMOption) (OOSMServiceClient, error) {
	svc, err := NewOutOfStoreMessageService(opts...)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	c, err := newClientFromService(ctx, s, svc)
	if err != nil {
		return nil, fmt.Errorf("uanble to create client from server: %w", err)
	}

	return &oosmServiceClient{
		OOSMServiceClient: c,
		server:            s,
		service:           svc,
	}, nil
}

type oosmClient struct {
	outofstoremessagetypes.OutOfStoreMessageServiceClient

	l  *grpcutil.BufListener
	cc *grpc.ClientConn
}

func (c *oosmClient) Close() error {
	err := c.cc.Close()
	_ = c.l.Close()
	return err
}

func newClientFromService(ctx context.Context, s *grpc.Server, svc OOSMService, opts ...grpc.DialOption) (OOSMServiceClient, error) {
	bl := grpcutil.NewBufListener(ClientBufferSize)
	cc, err := bl.NewClientConn(ctx, opts...)
	if err != nil {
		return nil, err
	}

	outofstoremessagetypes.RegisterOutOfStoreMessageServiceServer(s, svc)
	go func() {
		// we dont need to log the error
		_ = s.Serve(bl)
	}()

	return &oosmClient{
		OutOfStoreMessageServiceClient: outofstoremessagetypes.NewOutOfStoreMessageServiceClient(cc),
		cc:                             cc,
		l:                              bl,
	}, nil
}

func NewOutOfStoreMessageService(opts ...OOSMOption) (OOSMService, error) {
	svc := &oosmService{}

	withDefaultOpts := make([]OOSMOption, len(opts))
	copy(withDefaultOpts, opts)
	withDefaultOpts = append(withDefaultOpts, WithFallbackDefaults)
	for _, opt := range withDefaultOpts {
		if err := opt(svc); err != nil {
			return nil, err
		}
	}

	return svc, nil
}

func (s *oosmService) Close() error {
	return nil
}

func (s *oosmService) Status() (Status, error) {
	return Status{}, nil
}

func (s *oosmService) IpfsCoreAPI() coreiface.CoreAPI {
	return nil
}

func (s *oosmService) OutOfStoreReceive(ctx context.Context, request *protocoltypes.OutOfStoreReceive_Request) (*protocoltypes.OutOfStoreReceive_Reply, error) {
	outOfStoreMessage, group, clearPayload, alreadyDecrypted, err := s.secretStore.OpenOutOfStoreMessage(ctx, request.Payload)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoDecrypt.Wrap(err)
	}

	return &protocoltypes.OutOfStoreReceive_Reply{
		Message:         outOfStoreMessage,
		Cleartext:       clearPayload,
		GroupPublicKey:  group.PublicKey,
		AlreadyReceived: alreadyDecrypted,
	}, nil
}

// FallBackOption is a structure that permit to fallback to a default option if the option is not set.
type FallBackOption struct {
	fallback func(s *oosmService) bool
	opt      OOSMOption
}

// WithLogger set the given logger.
var WithLogger = func(l *zap.Logger) OOSMOption {
	return func(s *oosmService) error {
		s.logger = l
		return nil
	}
}

// WithDefaultLogger init a noop logger.
var WithDefaultLogger OOSMOption = func(s *oosmService) error {
	s.logger = zap.NewNop()
	return nil
}

var fallbackLogger = FallBackOption{
	fallback: func(s *oosmService) bool { return s.logger == nil },
	opt:      WithDefaultLogger,
}

// WithFallbackLogger set the logger if no logger is set.
var WithFallbackLogger OOSMOption = func(s *oosmService) error {
	if fallbackLogger.fallback(s) {
		return fallbackLogger.opt(s)
	}
	return nil
}

// WithRootDatastore set the root datastore.
var WithRootDatastore = func(ds ds.Datastore) OOSMOption {
	return func(s *oosmService) error {
		s.rootDatastore = ds
		return nil
	}
}

// WithDefaultRootDatastore init a in-memory datastore.
var WithDefaultRootDatastore OOSMOption = func(s *oosmService) error {
	s.rootDatastore = ds_sync.MutexWrap(ds.NewMapDatastore())
	return nil
}

var fallbackRootDatastore = FallBackOption{
	fallback: func(s *oosmService) bool { return s.rootDatastore == nil },
	opt:      WithDefaultRootDatastore,
}

// WithFallbackRootDatastore set the root datastore if no root datastore is set.
var WithFallbackRootDatastore OOSMOption = func(s *oosmService) error {
	if fallbackRootDatastore.fallback(s) {
		return fallbackRootDatastore.opt(s)
	}
	return nil
}

// WithSecretStore set the secret store.
var WithSecretStore = func(ss secretstore.SecretStore) OOSMOption {
	return func(s *oosmService) error {
		s.secretStore = ss
		return nil
	}
}

// WithDefaultSecretStore init a new secret store.
// Call WithRootDatastore before this option if you want to use your datastore.
// Call WithLogger before this option if you want to use your logger.
var WithDefaultSecretStore OOSMOption = func(s *oosmService) error {
	// dependency
	if err := WithFallbackRootDatastore(s); err != nil {
		return err
	}
	if err := WithFallbackLogger(s); err != nil {
		return err
	}

	var err error
	s.secretStore, err = secretstore.NewSecretStore(s.rootDatastore, &secretstore.NewSecretStoreOptions{
		Logger: s.logger,
	})
	return err
}

var fallbackSecretStore = FallBackOption{
	fallback: func(s *oosmService) bool { return s.secretStore == nil },
	opt:      WithDefaultSecretStore,
}

// WithFallbackSecretStore set the secret store if no secret store is set.
// Call WithRootDatastore before this option if you want to use your datastore if a new secret store is created.
// Call WithLogger before this option if you want to use your logger if a new secret store is created.
var WithFallbackSecretStore OOSMOption = func(s *oosmService) error {
	if fallbackSecretStore.fallback(s) {
		return fallbackSecretStore.opt(s)
	}
	return nil
}

var defaults = []FallBackOption{
	fallbackLogger,
	fallbackRootDatastore,
	fallbackSecretStore,
}

// WithFallbackDefaults set the default options if no option is set.
var WithFallbackDefaults OOSMOption = func(s *oosmService) error {
	for _, def := range defaults {
		if !def.fallback(s) {
			continue
		}
		if err := def.opt(s); err != nil {
			return err
		}
	}
	return nil
}
