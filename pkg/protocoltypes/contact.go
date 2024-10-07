package protocoltypes

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"

	"berty.tech/weshnet/v2/pkg/errcode"
)

const RendezvousSeedLength = 32

type ShareableContactOptions uint64

const (
	shareableContactOptionsUndefined = iota
	ShareableContactOptionsAllowMissingRDVSeed
	ShareableContactOptionsAllowMissingPK
)

var _ = shareableContactOptionsUndefined

func (m *ShareableContact) CheckFormat(options ...ShareableContactOptions) error {
	var (
		optionMissingPKAllowed      = false
		optionMissingRDVSeedAllowed = false
	)

	for _, o := range options {
		if o == ShareableContactOptionsAllowMissingPK {
			optionMissingPKAllowed = true
		}

		if o == ShareableContactOptionsAllowMissingRDVSeed {
			optionMissingRDVSeedAllowed = true
		}
	}

	if l := len(m.PublicRendezvousSeed); l != RendezvousSeedLength {
		if !(l == 0 && optionMissingRDVSeedAllowed) {
			return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("rendezvous seed length should not be %d", l))
		}
	}

	if l := len(m.Pk); l == 0 && !optionMissingPKAllowed {
		return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("contact public key is missing"))
	}

	if l := len(m.Pk); l != 0 {
		_, err := crypto.UnmarshalEd25519PublicKey(m.Pk)
		if err != nil {
			return errcode.ErrCode_ErrDeserialization.Wrap(err)
		}
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
	pk, err := crypto.UnmarshalEd25519PublicKey(m.Pk)
	if err != nil {
		return nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	return pk, nil
}
