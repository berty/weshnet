package bertyprotocol

import (
	"github.com/libp2p/go-libp2p-core/crypto"

	"berty.tech/berty/go/pkg/errcode"
)

const RendezvousSeedLength = 32

func (m *ShareableContact) CheckFormat() error {
	if len(m.PublicRendezvousSeed) != RendezvousSeedLength {
		return errcode.ErrInvalidInput
	}

	_, err := crypto.UnmarshalEd25519PublicKey(m.PK)
	if err != nil {
		return errcode.ErrDeserialization.Wrap(err)
	}

	return nil
}
func (m *ShareableContact) IsSamePK(otherPK crypto.PubKey) bool {
	pk, err := m.GetPubKey()
	if err != nil {
		return false
	}

	return otherPK.Equals(pk)
}

func (m *ShareableContact) GetPubKey() (crypto.PubKey, error) {
	pk, err := crypto.UnmarshalEd25519PublicKey(m.PK)
	if err != nil {
		return nil, errcode.ErrDeserialization.Wrap(err)
	}

	return pk, nil
}
