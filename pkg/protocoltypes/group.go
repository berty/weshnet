package protocoltypes

import (
	crand "crypto/rand"
	"encoding/hex"
	"io"

	"github.com/libp2p/go-libp2p/core/crypto"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/sha3"

	"berty.tech/weshnet/v2/pkg/cryptoutil"
	"berty.tech/weshnet/v2/pkg/errcode"
)

func (m *Group) GetSigningPrivKey() (crypto.PrivKey, error) {
	if len(m.Secret) == 0 {
		return nil, errcode.ErrCode_ErrMissingInput
	}

	edSK := ed25519.NewKeyFromSeed(m.Secret)

	sk, _, err := crypto.KeyPairFromStdKey(&edSK)
	if err != nil {
		return nil, err
	}

	return sk, nil
}

func (m *Group) GetPubKey() (crypto.PubKey, error) {
	return crypto.UnmarshalEd25519PublicKey(m.PublicKey)
}

func (m *Group) GetSigningPubKey() (crypto.PubKey, error) {
	if len(m.SignPub) != 0 {
		return crypto.UnmarshalEd25519PublicKey(m.SignPub)
	}

	sk, err := m.GetSigningPrivKey()
	if err != nil {
		return nil, err
	}

	return sk.GetPublic(), nil
}

func (m *Group) IsValid() error {
	pk, err := m.GetPubKey()
	if err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	ok, err := pk.Verify(m.Secret, m.SecretSig)
	if err != nil {
		return errcode.ErrCode_ErrCryptoSignatureVerification.Wrap(err)
	}

	if !ok {
		return errcode.ErrCode_ErrCryptoSignatureVerification
	}

	return nil
}

// GroupIDAsString returns the group pub key as a string
func (m *Group) GroupIDAsString() string {
	return hex.EncodeToString(m.PublicKey)
}

func (m *Group) Copy() *Group {
	return &Group{
		PublicKey: m.PublicKey,
		Secret:    m.Secret,
		SecretSig: m.SecretSig,
		GroupType: m.GroupType,
		SignPub:   m.SignPub,
	}
}

const CurrentGroupVersion = 1

// NewGroupMultiMember creates a new Group object and an invitation to be used by
// the first member of the group
func NewGroupMultiMember() (*Group, crypto.PrivKey, error) {
	priv, pub, err := crypto.GenerateEd25519Key(crand.Reader)
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	pubBytes, err := pub.Raw()
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	signing, _, err := crypto.GenerateEd25519Key(crand.Reader)
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	signingBytes, err := cryptoutil.SeedFromEd25519PrivateKey(signing)
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	skSig, err := priv.Sign(signingBytes)
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrCryptoSignature.Wrap(err)
	}

	group := &Group{
		PublicKey: pubBytes,
		Secret:    signingBytes,
		SecretSig: skSig,
		GroupType: GroupType_GroupTypeMultiMember,
	}

	updateKey, err := group.GetLinkKeyArray()
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	linkKeySig, err := priv.Sign(updateKey[:])
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrCryptoSignature.Wrap(err)
	}

	group.LinkKeySig = linkKeySig

	return group, priv, nil
}

func ComputeLinkKey(publicKey, secret []byte) (*[cryptoutil.KeySize]byte, error) {
	arr := [cryptoutil.KeySize]byte{}

	kdf := hkdf.New(sha3.New256, secret, nil, publicKey)
	if _, err := io.ReadFull(kdf, arr[:]); err != nil {
		return nil, errcode.ErrCode_ErrStreamRead.Wrap(err)
	}

	return &arr, nil
}

func (m *Group) GetLinkKeyArray() (*[cryptoutil.KeySize]byte, error) {
	if len(m.GetLinkKey()) == cryptoutil.KeySize {
		arr := [cryptoutil.KeySize]byte{}

		for i, c := range m.GetLinkKey() {
			arr[i] = c
		}

		return &arr, nil
	}

	return ComputeLinkKey(m.GetPublicKey(), m.GetSecret())
}

func (m *Group) GetSharedSecret() *[cryptoutil.KeySize]byte {
	sharedSecret := [cryptoutil.KeySize]byte{}
	copy(sharedSecret[:], m.GetSecret())

	return &sharedSecret
}
