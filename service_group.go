package weshnet

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"
	"go.uber.org/zap"

	"berty.tech/go-orbit-db/iface"
	"berty.tech/weshnet/v2/pkg/errcode"
	"berty.tech/weshnet/v2/pkg/protocoltypes"
	"berty.tech/weshnet/v2/pkg/secretstore"
)

func (s *service) getContactGroup(key crypto.PubKey) (*protocoltypes.Group, error) {
	group, err := s.secretStore.GetGroupForContact(key)
	if err != nil {
		return nil, errcode.ErrCode_ErrOrbitDBOpen.Wrap(err)
	}

	return group, nil
}

func (s *service) getGroupForPK(ctx context.Context, pk crypto.PubKey) (*protocoltypes.Group, error) {
	group, err := s.secretStore.FetchGroupByPublicKey(ctx, pk)
	if err == nil {
		return group, nil
	} else if !errcode.Is(err, errcode.ErrCode_ErrMissingMapKey) {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	accountGroup := s.getAccountGroup()
	if accountGroup == nil {
		return nil, errcode.ErrCode_ErrGroupMissing
	}

	if err = reindexGroupDatastore(ctx, s.secretStore, accountGroup.metadataStore); err != nil {
		return nil, errcode.ErrCode_TODO.Wrap(err)
	}

	group, err = s.secretStore.FetchGroupByPublicKey(ctx, pk)
	if err == nil {
		return group, nil
	} else if errcode.Is(err, errcode.ErrCode_ErrMissingMapKey) {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("unknown group specified"))
	}

	return nil, errcode.ErrCode_ErrInternal.Wrap(err)
}

func (s *service) deactivateGroup(pk crypto.PubKey) error {
	id, err := pk.Raw()
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	cg, err := s.GetContextGroupForID(id)
	if err != nil || cg == nil {
		// @FIXME(gfanton): should return an error code
		return nil
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	err = cg.Close()
	if err != nil {
		s.logger.Error("unable to close group context", zap.Error(err))
	}

	delete(s.openedGroups, string(id))

	if cg.group.GroupType == protocoltypes.GroupType_GroupTypeAccount {
		s.accountGroupCtx = nil
	}

	return nil
}

func (s *service) activateGroup(ctx context.Context, pk crypto.PubKey, localOnly bool) error {
	id, err := pk.Raw()
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	_, err = s.GetContextGroupForID(id)
	if err != nil && err != errcode.ErrCode_ErrGroupUnknown {
		return err
	}

	g, err := s.getGroupForPK(ctx, pk)
	if err != nil {
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	// @WIP(gfanton): do we need to use contactPK
	var contactPK crypto.PubKey
	switch g.GroupType {
	case protocoltypes.GroupType_GroupTypeMultiMember:
		// nothing to get here, simply continue, open and activate the group

	case protocoltypes.GroupType_GroupTypeContact:
		if s.accountGroupCtx == nil {
			return errcode.ErrCode_ErrGroupActivate.Wrap(fmt.Errorf("accountGroupCtx is deactivated"))
		}

		contact := s.accountGroupCtx.metadataStore.GetContactFromGroupPK(id)
		if contact != nil {
			contactPK, err = contact.GetPubKey()
			if err != nil {
				return errcode.ErrCode_TODO.Wrap(err)
			}
		}
	case protocoltypes.GroupType_GroupTypeAccount:
		localOnly = true
		if s.accountGroupCtx, err = s.odb.openAccountGroup(ctx, &iface.CreateDBOptions{EventBus: s.accountEventBus, LocalOnly: &localOnly}, s.ipfsCoreAPI); err != nil {
			return err
		}
		s.openedGroups[string(id)] = s.accountGroupCtx

		// reinitialize contactRequestsManager
		if s.contactRequestsManager != nil {
			s.contactRequestsManager.close()

			if s.contactRequestsManager, err = newContactRequestsManager(s.swiper, s.accountGroupCtx.metadataStore, s.ipfsCoreAPI, s.logger); err != nil {
				return errcode.ErrCode_TODO.Wrap(err)
			}
		}
		return nil
	default:
		return errcode.ErrCode_ErrInternal.Wrap(fmt.Errorf("unknown group type"))
	}

	dbOpts := &iface.CreateDBOptions{LocalOnly: &localOnly}
	gc, err := s.odb.OpenGroup(ctx, g, dbOpts)
	if err != nil {
		return errcode.ErrCode_ErrGroupOpen.Wrap(err)
	}

	if err = gc.ActivateGroupContext(contactPK); err != nil {
		gc.Close()
		return errcode.ErrCode_ErrGroupActivate.Wrap(err)
	}

	s.openedGroups[string(id)] = gc

	gc.TagGroupContextPeers(s.ipfsCoreAPI, 42)
	return nil
}

func (s *service) GetContextGroupForID(id []byte) (*GroupContext, error) {
	if len(id) == 0 {
		return nil, errcode.ErrCode_ErrInternal.Wrap(fmt.Errorf("no group id provided"))
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	cg, ok := s.openedGroups[string(id)]

	if ok {
		return cg, nil
	}

	return nil, errcode.ErrCode_ErrGroupUnknown
}

func reindexGroupDatastore(ctx context.Context, secretStore secretstore.SecretStore, m *MetadataStore) error {
	if secretStore == nil {
		return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("missing device keystore"))
	}

	for _, g := range m.ListMultiMemberGroups() {
		if err := secretStore.PutGroup(ctx, g); err != nil {
			return errcode.ErrCode_ErrInternal.Wrap(err)
		}
	}

	for _, contact := range m.ListContactsByStatus(
		protocoltypes.ContactState_ContactStateToRequest,
		protocoltypes.ContactState_ContactStateReceived,
		protocoltypes.ContactState_ContactStateAdded,
		protocoltypes.ContactState_ContactStateRemoved,
		protocoltypes.ContactState_ContactStateDiscarded,
		protocoltypes.ContactState_ContactStateBlocked,
	) {
		cPK, err := contact.GetPubKey()
		if err != nil {
			return errcode.ErrCode_TODO.Wrap(err)
		}

		group, err := secretStore.GetGroupForContact(cPK)
		if err != nil {
			return errcode.ErrCode_ErrInternal.Wrap(err)
		}

		if err := secretStore.PutGroup(ctx, group); err != nil {
			return errcode.ErrCode_ErrInternal.Wrap(err)
		}
	}

	return nil
}

func (s *service) getAccountGroup() *GroupContext {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.accountGroupCtx
}
