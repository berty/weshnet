package secretstore

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p/core/crypto"
	crypto_pb "github.com/libp2p/go-libp2p/core/crypto/pb"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/sha3"

	"berty.tech/weshnet/v2/pkg/cryptoutil"
	"berty.tech/weshnet/v2/pkg/errcode"
	"berty.tech/weshnet/v2/pkg/protocoltypes"
)

// getEd25519PrivateKeyFromLibP2PFormattedBytes transforms an exported LibP2P
// private key into a crypto.PrivKey instance, ensuring it is an ed25519 key
func getEd25519PrivateKeyFromLibP2PFormattedBytes(rawKeyBytes []byte) (crypto.PrivKey, error) {
	if len(rawKeyBytes) == 0 {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("missing key data"))
	}

	privateKey, err := crypto.UnmarshalPrivateKey(rawKeyBytes)
	if err != nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	if privateKey.Type() != crypto_pb.KeyType_Ed25519 {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("invalid key format"))
	}

	return privateKey, nil
}

// getKeysForGroupOfContact returns derived private keys for contact group
// using a private key via two accounts account keys (via an ECDH).
func getKeysForGroupOfContact(contactPairPrivateKey crypto.PrivKey) (crypto.PrivKey, crypto.PrivKey, error) {
	// Salt length must be equal to hash length (64 bytes for sha256)
	hash := sha256.New

	contactPairPrivateKeyBytes, err := contactPairPrivateKey.Raw()
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	// Generate Pseudo Random Key using contactPairPrivateKeyBytes as IKM and salt
	prk := hkdf.Extract(hash, contactPairPrivateKeyBytes, nil)
	if len(prk) == 0 {
		return nil, nil, errcode.ErrCode_ErrInternal.Wrap(fmt.Errorf("unable to instantiate pseudo random key"))
	}

	// Expand using extracted prk and groupID as info (kind of namespace)
	kdf := hkdf.Expand(hash, prk, nil)

	// Generate next KDF and message keys
	groupSeed, err := io.ReadAll(io.LimitReader(kdf, 32))
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	groupSecretSeed, err := io.ReadAll(io.LimitReader(kdf, 32))
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	stdGroupPrivateKey := ed25519.NewKeyFromSeed(groupSeed)
	groupPrivateKey, _, err := crypto.KeyPairFromStdKey(&stdGroupPrivateKey)
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	stdGroupSecretPrivateKey := ed25519.NewKeyFromSeed(groupSecretSeed)
	groupSecretPrivateKey, _, err := crypto.KeyPairFromStdKey(&stdGroupSecretPrivateKey)
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	return groupPrivateKey, groupSecretPrivateKey, nil
}

// getGroupForContact returns a protocoltypes.Group instance for a contact,
// using a private key via two accounts account keys (via an ECDH)
func getGroupForContact(contactPairPrivateKey crypto.PrivKey) (*protocoltypes.Group, error) {
	groupPrivateKey, groupSecretPrivateKey, err := getKeysForGroupOfContact(contactPairPrivateKey)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	pubBytes, err := groupPrivateKey.GetPublic().Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	signingBytes, err := cryptoutil.SeedFromEd25519PrivateKey(groupSecretPrivateKey)
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	return &protocoltypes.Group{
		PublicKey: pubBytes,
		Secret:    signingBytes,
		SecretSig: nil,
		GroupType: protocoltypes.GroupType_GroupTypeContact,
	}, nil
}

// getGroupOutOfStoreSecret retrieves the out of store group secret
func getGroupOutOfStoreSecret(m *protocoltypes.Group) ([]byte, error) {
	if len(m.GetSecret()) == 0 {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("no secret known for group"))
	}

	arr := [cryptoutil.KeySize]byte{}

	kdf := hkdf.New(sha3.New256, m.GetSecret(), nil, []byte(namespaceOutOfStoreSecret))
	if _, err := io.ReadFull(kdf, arr[:]); err != nil {
		return nil, errcode.ErrCode_ErrStreamRead.Wrap(err)
	}

	return arr[:], nil
}

// createOutOfStoreGroupReference creates a hash used to identify an out of
// store (e.g. push notification) message origin
func createOutOfStoreGroupReference(m *protocoltypes.Group, sender []byte, counter uint64) ([]byte, error) {
	secret, err := getGroupOutOfStoreSecret(m)
	if err != nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	arr := [cryptoutil.KeySize]byte{}

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)

	kdf := hkdf.New(sha3.New256, secret, nil, append(sender, buf...))
	if _, err := io.ReadFull(kdf, arr[:]); err != nil {
		return nil, errcode.ErrCode_ErrStreamRead.Wrap(err)
	}

	return arr[:], nil
}
