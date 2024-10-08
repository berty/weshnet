package weshnet

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"

	"berty.tech/weshnet/v2/pkg/errcode"
	"berty.tech/weshnet/v2/pkg/protocoltypes"
	"berty.tech/weshnet/v2/pkg/tyber"
)

func (s *service) ContactAliasKeySend(ctx context.Context, req *protocoltypes.ContactAliasKeySend_Request) (_ *protocoltypes.ContactAliasKeySend_Reply, err error) {
	ctx, _, endSection := tyber.Section(ctx, s.logger, "Sending contact alias key")
	defer func() { endSection(err, "") }()

	g, err := s.GetContextGroupForID(req.GroupPk)
	if err != nil {
		return nil, errcode.ErrCode_ErrGroupMissing.Wrap(err)
	}

	if _, err := g.MetadataStore().ContactSendAliasKey(ctx); err != nil {
		return nil, errcode.ErrCode_ErrOrbitDBAppend.Wrap(err)
	}

	return &protocoltypes.ContactAliasKeySend_Reply{}, nil
}

func (s *service) ContactBlock(ctx context.Context, req *protocoltypes.ContactBlock_Request) (_ *protocoltypes.ContactBlock_Reply, err error) {
	ctx, _, endSection := tyber.Section(ctx, s.logger, "Blocking contact")
	defer func() { endSection(err, "") }()

	pk, err := crypto.UnmarshalEd25519PublicKey(req.ContactPk)
	if err != nil {
		return nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	if _, err := s.getAccountGroup().MetadataStore().ContactBlock(ctx, pk); err != nil {
		return nil, errcode.ErrCode_ErrOrbitDBAppend.Wrap(err)
	}

	return &protocoltypes.ContactBlock_Reply{}, nil
}

func (s *service) ContactUnblock(ctx context.Context, req *protocoltypes.ContactUnblock_Request) (_ *protocoltypes.ContactUnblock_Reply, err error) {
	ctx, _, endSection := tyber.Section(ctx, s.logger, "Unblocking contact")
	defer func() { endSection(err, "") }()

	pk, err := crypto.UnmarshalEd25519PublicKey(req.ContactPk)
	if err != nil {
		return nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	if _, err := s.getAccountGroup().MetadataStore().ContactUnblock(ctx, pk); err != nil {
		return nil, errcode.ErrCode_ErrOrbitDBAppend.Wrap(err)
	}

	return &protocoltypes.ContactUnblock_Reply{}, nil
}

func (s *service) RefreshContactRequest(ctx context.Context, req *protocoltypes.RefreshContactRequest_Request) (*protocoltypes.RefreshContactRequest_Reply, error) {
	if len(req.ContactPk) == 0 {
		return nil, errcode.ErrCode_ErrInternal
	}

	var cancel context.CancelFunc
	if req.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(req.Timeout)*time.Second)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	key := string(req.ContactPk)
	s.muRefreshprocess.Lock()
	if clfn, ok := s.refreshprocess[key]; ok {
		clfn() // close previous refresh method
	}
	s.refreshprocess[key] = cancel
	s.muRefreshprocess.Unlock()

	peers, err := s.swiper.RefreshContactRequest(ctx, req.ContactPk)
	if err != nil {
		return nil, fmt.Errorf("unable to refresh group: %w", err)
	}

	res := &protocoltypes.RefreshContactRequest_Reply{
		PeersFound: []*protocoltypes.RefreshContactRequest_Peer{},
	}
	for _, p := range peers {
		// check if we can connect to this peers
		if err := s.host.Connect(ctx, p); err != nil {
			continue
		}

		addrs := make([]string, len(p.Addrs))
		for i, addr := range p.Addrs {
			addrs[i] = addr.String()
		}

		res.PeersFound = append(res.PeersFound, &protocoltypes.RefreshContactRequest_Peer{
			Id:    p.ID.String(),
			Addrs: addrs,
		})
	}

	return res, nil
}
