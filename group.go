package weshnet

import (
	"github.com/libp2p/go-libp2p/core/crypto"

	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/protocoltypes"
)

const CurrentGroupVersion = 1

// NewGroupMultiMember creates a new Group object and an invitation to be used by
// the first member of the group
func NewGroupMultiMember() (*protocoltypes.Group, crypto.PrivKey, error) {
	return protocoltypes.NewGroupMultiMember()
}

func getAndFilterGroupAddDeviceChainKeyPayload(m *protocoltypes.GroupMetadata, localMemberPublicKey crypto.PubKey) (crypto.PubKey, []byte, error) {
	if m == nil || m.EventType != protocoltypes.EventTypeGroupDeviceChainKeyAdded {
		return nil, nil, errcode.ErrInvalidInput
	}

	s := &protocoltypes.GroupAddDeviceChainKey{}
	if err := s.Unmarshal(m.Payload); err != nil {
		return nil, nil, errcode.ErrDeserialization.Wrap(err)
	}

	senderDevicePubKey, err := crypto.UnmarshalEd25519PublicKey(s.DevicePK)
	if err != nil {
		return nil, nil, errcode.ErrDeserialization.Wrap(err)
	}

	destMemberPubKey, err := crypto.UnmarshalEd25519PublicKey(s.DestMemberPK)
	if err != nil {
		return nil, nil, errcode.ErrDeserialization.Wrap(err)
	}

	if !localMemberPublicKey.Equals(destMemberPubKey) {
		return nil, nil, errcode.ErrGroupSecretOtherDestMember
	}

	return senderDevicePubKey, s.Payload, nil
}
