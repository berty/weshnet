package weshnet

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/ipfs/go-cid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/logutil"
	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/pushtypes"
)

func (s *service) getPushClient(host string) (pushtypes.PushServiceClient, error) {
	s.muPushClients.Lock()
	defer s.muPushClients.Unlock()

	if cc, ok := s.pushClients[host]; ok {
		return pushtypes.NewPushServiceClient(cc), nil
	}

	var creds grpc.DialOption
	if s.grpcInsecure {
		creds = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		tlsconfig := credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})
		creds = grpc.WithTransportCredentials(tlsconfig)
	}

	// retry policies
	connectParams := grpc.WithConnectParams(grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay:  1.0 * time.Second,
			Multiplier: 1.5,
			Jitter:     0.2,
			MaxDelay:   60 * time.Second,
		},
		MinConnectTimeout: time.Second * 10,
	})

	cc, err := grpc.DialContext(s.ctx, host, creds, connectParams)
	if err != nil {
		return nil, err
	}
	s.pushClients[host] = cc

	// monitor push client state
	go monitorPushServer(s.ctx, cc, s.logger)

	return pushtypes.NewPushServiceClient(cc), err
}

func (s *service) OutOfStoreReceive(ctx context.Context, request *protocoltypes.OutOfStoreReceive_Request) (*protocoltypes.OutOfStoreReceive_Reply, error) {
	outOfStoreMessage, group, clearPayload, alreadyDecrypted, err := s.secretStore.OpenOutOfStoreMessage(ctx, request.Payload)
	if err != nil {
		return nil, errcode.ErrCryptoDecrypt.Wrap(err)
	}

	return &protocoltypes.OutOfStoreReceive_Reply{
		Message:         outOfStoreMessage,
		Cleartext:       clearPayload,
		GroupPublicKey:  group.PublicKey,
		AlreadyReceived: alreadyDecrypted,
	}, nil
}

func (s *service) OutOfStoreSeal(ctx context.Context, request *protocoltypes.OutOfStoreSeal_Request) (*protocoltypes.OutOfStoreSeal_Reply, error) {
	gc, err := s.GetContextGroupForID(request.GroupPublicKey)
	if err != nil {
		return nil, err
	}

	_, c, err := cid.CidFromBytes(request.CID)
	if err != nil {
		return nil, errcode.ErrInvalidInput.Wrap(err)
	}

	sealedMessageEnvelope, err := gc.messageStore.GetOutOfStoreMessageEnvelope(ctx, c)
	if err != nil {
		return nil, errcode.ErrInternal.Wrap(err)
	}

	sealedMessageEnvelopeBytes, err := sealedMessageEnvelope.Marshal()
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	return &protocoltypes.OutOfStoreSeal_Reply{
		Encrypted: sealedMessageEnvelopeBytes,
	}, nil
}

func (s *service) GetCurrentDevicePushConfig() (*protocoltypes.PushServiceReceiver, *protocoltypes.PushServer) {
	accountGroup := s.getAccountGroup()
	if accountGroup == nil {
		return nil, nil
	}

	currentToken := accountGroup.metadataStore.getCurrentDevicePushToken()
	currentServer := accountGroup.metadataStore.getCurrentDevicePushServer()

	return currentToken, currentServer
}

func monitorPushServer(ctx context.Context, cc *grpc.ClientConn, logger *zap.Logger) {
	currentState := cc.GetState()
	for cc.WaitForStateChange(ctx, currentState) {
		currentState = cc.GetState()
		logger.Debug("push grpc client state updated",
			logutil.PrivateString("target", cc.Target()),
			logutil.PrivateString("state", currentState.String()))
	}
}
