package weshnet

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"

	"berty.tech/weshnet/v2/pkg/errcode"
	"berty.tech/weshnet/v2/pkg/protocoltypes"
	"berty.tech/weshnet/v2/pkg/tyber"
)

// MultiMemberGroupCreate creates a new MultiMember group
func (s *service) MultiMemberGroupCreate(ctx context.Context, _ *protocoltypes.MultiMemberGroupCreate_Request) (_ *protocoltypes.MultiMemberGroupCreate_Reply, err error) {
	ctx, _, endSection := tyber.Section(ctx, s.logger, "Creating MultiMember group")
	defer func() { endSection(err, "") }()

	group, groupPrivateKey, err := NewGroupMultiMember()
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	accountGroup := s.getAccountGroup()
	if accountGroup == nil {
		return nil, errcode.ErrCode_ErrGroupMissing
	}

	_, err = accountGroup.MetadataStore().GroupJoin(ctx, group)
	if err != nil {
		return nil, errcode.ErrCode_ErrOrbitDBAppend.Wrap(err)
	}

	if err := s.secretStore.PutGroup(ctx, group); err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	err = s.activateGroup(ctx, groupPrivateKey.GetPublic(), false)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(fmt.Errorf("unable to activate group: %w", err))
	}

	cg, err := s.GetContextGroupForID(group.PublicKey)
	if err != nil {
		return nil, errcode.ErrCode_ErrOrbitDBAppend.Wrap(err)
	}

	_, err = cg.MetadataStore().ClaimGroupOwnership(ctx, groupPrivateKey)
	if err != nil {
		return nil, errcode.ErrCode_ErrOrbitDBAppend.Wrap(err)
	}

	return &protocoltypes.MultiMemberGroupCreate_Reply{
		GroupPk: group.PublicKey,
	}, nil
}

// MultiMemberGroupJoin joins an existing MultiMember group using an invitation
func (s *service) MultiMemberGroupJoin(ctx context.Context, req *protocoltypes.MultiMemberGroupJoin_Request) (_ *protocoltypes.MultiMemberGroupJoin_Reply, err error) {
	ctx, _, endSection := tyber.Section(ctx, s.logger, "Joining MultiMember group")
	defer func() { endSection(err, "") }()

	accountGroup := s.getAccountGroup()
	if accountGroup == nil {
		return nil, errcode.ErrCode_ErrGroupMissing
	}

	if _, err := accountGroup.MetadataStore().GroupJoin(ctx, req.Group); err != nil {
		return nil, errcode.ErrCode_ErrOrbitDBAppend.Wrap(err)
	}

	return &protocoltypes.MultiMemberGroupJoin_Reply{}, nil
}

// MultiMemberGroupLeave leaves a previously joined MultiMember group
func (s *service) MultiMemberGroupLeave(ctx context.Context, req *protocoltypes.MultiMemberGroupLeave_Request) (_ *protocoltypes.MultiMemberGroupLeave_Reply, err error) {
	ctx, _, endSection := tyber.Section(ctx, s.logger, "Leaving MultiMember group")
	defer func() { endSection(err, "") }()

	pk, err := crypto.UnmarshalEd25519PublicKey(req.GroupPk)
	if err != nil {
		return nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	accountGroup := s.getAccountGroup()
	if accountGroup == nil {
		return nil, errcode.ErrCode_ErrGroupMissing
	}

	_, err = accountGroup.MetadataStore().GroupLeave(ctx, pk)
	if err != nil {
		return nil, errcode.ErrCode_ErrOrbitDBAppend.Wrap(err)
	}

	if err := s.deactivateGroup(pk); err != nil {
		return nil, errcode.ErrCode_ErrOrbitDBAppend.Wrap(err)
	}

	return &protocoltypes.MultiMemberGroupLeave_Reply{}, nil
}

// MultiMemberGroupAliasResolverDisclose sends an deviceKeystore identity proof to the group members
func (s *service) MultiMemberGroupAliasResolverDisclose(ctx context.Context, req *protocoltypes.MultiMemberGroupAliasResolverDisclose_Request) (*protocoltypes.MultiMemberGroupAliasResolverDisclose_Reply, error) {
	cg, err := s.GetContextGroupForID(req.GroupPk)
	if err != nil {
		return nil, errcode.ErrCode_ErrGroupMemberUnknownGroupID.Wrap(err)
	}

	_, err = cg.MetadataStore().SendAliasProof(ctx)
	if err != nil {
		return nil, errcode.ErrCode_ErrOrbitDBAppend.Wrap(err)
	}

	return &protocoltypes.MultiMemberGroupAliasResolverDisclose_Reply{}, nil
}

// MultiMemberGroupAdminRoleGrant grants admin role to another member of the group
func (s *service) MultiMemberGroupAdminRoleGrant(context.Context, *protocoltypes.MultiMemberGroupAdminRoleGrant_Request) (*protocoltypes.MultiMemberGroupAdminRoleGrant_Reply, error) {
	return nil, errcode.ErrCode_ErrNotImplemented
}

// MultiMemberGroupInvitationCreate creates a group invitation
func (s *service) MultiMemberGroupInvitationCreate(_ context.Context, req *protocoltypes.MultiMemberGroupInvitationCreate_Request) (*protocoltypes.MultiMemberGroupInvitationCreate_Reply, error) {
	cg, err := s.GetContextGroupForID(req.GroupPk)
	if err != nil {
		return nil, errcode.ErrCode_ErrGroupMemberUnknownGroupID.Wrap(err)
	}

	return &protocoltypes.MultiMemberGroupInvitationCreate_Reply{
		Group: cg.Group(),
	}, nil
}
