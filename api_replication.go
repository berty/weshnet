package weshnet

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/grpcutil"
	"berty.tech/weshnet/pkg/logutil"
	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/replicationtypes"
	"berty.tech/weshnet/pkg/tyber"
)

func FilterGroupForReplication(m *protocoltypes.Group) (*protocoltypes.Group, error) {
	groupSigPK, err := m.GetSigningPubKey()
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	groupSigPKBytes, err := groupSigPK.Raw()
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	linkKey, err := m.GetLinkKeyArray()
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	return &protocoltypes.Group{
		PublicKey:  m.PublicKey,
		SignPub:    groupSigPKBytes,
		LinkKey:    linkKey[:],
		LinkKeySig: m.LinkKeySig,
	}, nil
}

func (s *service) ReplicationServiceRegisterGroup(ctx context.Context, request *protocoltypes.ReplicationServiceRegisterGroup_Request) (_ *protocoltypes.ReplicationServiceRegisterGroup_Reply, err error) {
	ctx, _, endSection := tyber.Section(ctx, s.logger, "Registering replication service for group")
	defer func() { endSection(err, "") }()

	if request.GroupPK == nil {
		return nil, errcode.ErrInvalidInput.Wrap(fmt.Errorf("invalid GroupPK"))
	}

	if request.Token == "" {
		return nil, errcode.ErrInvalidInput.Wrap(fmt.Errorf("invalid token"))
	}

	if request.ReplicationServer == "" {
		return nil, errcode.ErrInvalidInput.Wrap(fmt.Errorf("invalid replication server"))
	}

	gc, err := s.GetContextGroupForID(request.GroupPK)
	if err != nil {
		return nil, errcode.ErrInvalidInput.Wrap(err)
	}

	replGroup, err := FilterGroupForReplication(gc.group)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	accountGroup := s.getAccountGroup()
	if accountGroup == nil {
		return nil, errcode.ErrGroupMissing
	}

	gopts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(grpcutil.NewUnsecureSimpleAuthAccess("bearer", request.Token)),
	}

	if s.grpcInsecure {
		gopts = append(gopts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tlsconfig := credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})
		gopts = append(gopts, grpc.WithTransportCredentials(tlsconfig))
	}

	cc, err := grpc.DialContext(context.Background(), request.ReplicationServer, gopts...)
	if err != nil {
		return nil, errcode.ErrStreamWrite.Wrap(err)
	}

	client := replicationtypes.NewReplicationServiceClient(cc)

	if _, err = client.ReplicateGroup(ctx, &replicationtypes.ReplicationServiceReplicateGroup_Request{
		Group: replGroup,
	}); err != nil {
		return nil, errcode.ErrServiceReplicationServer.Wrap(err)
	}

	s.logger.Info("group will be replicated", logutil.PrivateString("public-key", base64.RawURLEncoding.EncodeToString(request.GroupPK)))

	if _, err := gc.metadataStore.SendGroupReplicating(ctx, request.AuthenticationURL, request.ReplicationServer); err != nil {
		s.logger.Error("error while notifying group about replication", zap.Error(err))
	}

	return &protocoltypes.ReplicationServiceRegisterGroup_Reply{}, nil
}
