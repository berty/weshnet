package weshnet

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/ipfs/go-cid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/tyber"
)

func (s *service) AppMetadataSend(ctx context.Context, req *protocoltypes.AppMetadataSend_Request) (_ *protocoltypes.AppMetadataSend_Reply, err error) {
	ctx, _, endSection := tyber.Section(ctx, s.logger, fmt.Sprintf("Sending app metadata to group %s", base64.RawURLEncoding.EncodeToString(req.GroupPk)))
	defer func() { endSection(err, "") }()

	gc, err := s.GetContextGroupForID(req.GroupPk)
	if err != nil {
		return nil, errcode.ErrCode_ErrGroupMissing.Wrap(err)
	}
	tyberLogGroupContext(ctx, s.logger, gc)

	op, err := gc.MetadataStore().SendAppMetadata(ctx, req.Payload)
	if err != nil {
		return nil, errcode.ErrCode_ErrOrbitDBAppend.Wrap(err)
	}

	return &protocoltypes.AppMetadataSend_Reply{Cid: op.GetEntry().GetHash().Bytes()}, nil
}

func (s *service) AppMessageSend(ctx context.Context, req *protocoltypes.AppMessageSend_Request) (_ *protocoltypes.AppMessageSend_Reply, err error) {
	ctx, _, endSection := tyber.Section(ctx, s.logger, fmt.Sprintf("Sending message to group %s", base64.RawURLEncoding.EncodeToString(req.GroupPk)))
	defer func() { endSection(err, "") }()

	gc, err := s.GetContextGroupForID(req.GroupPk)
	if err != nil {
		return nil, errcode.ErrCode_ErrGroupMissing.Wrap(err)
	}
	tyberLogGroupContext(ctx, s.logger, gc)

	op, err := gc.MessageStore().AddMessage(ctx, req.Payload)
	if err != nil {
		return nil, errcode.ErrCode_ErrOrbitDBAppend.Wrap(err)
	}

	return &protocoltypes.AppMessageSend_Reply{Cid: op.GetEntry().GetHash().Bytes()}, nil
}

// OutOfStoreReceive parses a payload received outside a synchronized store
func (s *service) OutOfStoreReceive(ctx context.Context, request *protocoltypes.OutOfStoreReceive_Request) (*protocoltypes.OutOfStoreReceive_Reply, error) {
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

// OutOfStoreSeal creates a payload of a message present in store to be sent outside a synchronized store
func (s *service) OutOfStoreSeal(ctx context.Context, request *protocoltypes.OutOfStoreSeal_Request) (*protocoltypes.OutOfStoreSeal_Reply, error) {
	gc, err := s.GetContextGroupForID(request.GroupPublicKey)
	if err != nil {
		return nil, err
	}

	_, c, err := cid.CidFromBytes(request.Cid)
	if err != nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	sealedMessageEnvelope, err := gc.messageStore.GetOutOfStoreMessageEnvelope(ctx, c)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	sealedMessageEnvelopeBytes, err := proto.Marshal(sealedMessageEnvelope)
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	return &protocoltypes.OutOfStoreSeal_Reply{
		Encrypted: sealedMessageEnvelopeBytes,
	}, nil
}

func tyberLogGroupContext(ctx context.Context, logger *zap.Logger, gc *GroupContext) {
	memberPK, err := gc.MemberPubKey().Raw()
	if err != nil {
		memberPK = []byte{}
	}

	logger.Debug("Got group context", tyber.FormatStepLogFields(ctx, []tyber.Detail{
		{Name: "GroupType", Description: gc.Group().GetGroupType().String()},
		{Name: "GroupPK", Description: base64.RawURLEncoding.EncodeToString(gc.Group().PublicKey)},
		{Name: "MemberPK", Description: base64.RawURLEncoding.EncodeToString(memberPK)},
	})...)
}
