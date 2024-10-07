package weshnet

import (
	"github.com/libp2p/go-libp2p/core/crypto"
	"google.golang.org/protobuf/proto"

	"berty.tech/weshnet/v2/pkg/errcode"
	"berty.tech/weshnet/v2/pkg/protocoltypes"
)

const CurrentGroupVersion = 1

// NewGroupMultiMember creates a new Group object and an invitation to be used by
// the first member of the group
func NewGroupMultiMember() (*protocoltypes.Group, crypto.PrivKey, error) {
	return protocoltypes.NewGroupMultiMember()
}

func getAndFilterGroupDeviceChainKeyAddedPayload(m *protocoltypes.GroupMetadata, localMemberPublicKey crypto.PubKey) (crypto.PubKey, []byte, error) {
	if m == nil || m.EventType != protocoltypes.EventType_EventTypeGroupDeviceChainKeyAdded {
		return nil, nil, errcode.ErrCode_ErrInvalidInput
	}

	s := &protocoltypes.GroupDeviceChainKeyAdded{}
	if err := proto.Unmarshal(m.Payload, s); err != nil {
		return nil, nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	senderDevicePubKey, err := crypto.UnmarshalEd25519PublicKey(s.DevicePk)
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	destMemberPubKey, err := crypto.UnmarshalEd25519PublicKey(s.DestMemberPk)
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	if !localMemberPublicKey.Equals(destMemberPubKey) {
		return nil, nil, errcode.ErrCode_ErrGroupSecretOtherDestMember
	}

	return senderDevicePubKey, s.Payload, nil
}
