package secretstore

import (
	"crypto/ed25519"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"

	"github.com/aead/ecdh"
	keystore "github.com/ipfs/go-ipfs-keystore"
	"github.com/libp2p/go-libp2p/core/crypto"
	"go.uber.org/zap"

	"berty.tech/weshnet/v2/pkg/cryptoutil"
	"berty.tech/weshnet/v2/pkg/errcode"
	"berty.tech/weshnet/v2/pkg/protocoltypes"
)

const (
	keyAccount      = "accountSK"
	keyAccountProof = "accountProofSK"
	keyDevice       = "deviceSK"
	keyMemberDevice = "memberDeviceSK"
	keyMember       = "memberSK"
	keyContactGroup = "contactGroupSK"
)

// deviceKeystore is a wrapper around a keystore.Keystore object.
// It contains methods to manipulate member and device keys.
type deviceKeystore struct {
	keystore keystore.Keystore
	mu       sync.Mutex
	logger   *zap.Logger
}

// newDeviceKeystore instantiate a new device keystore
func newDeviceKeystore(ks keystore.Keystore, logger *zap.Logger) *deviceKeystore {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &deviceKeystore{
		keystore: ks,
		logger:   logger,
	}
}

// getAccountPrivateKey returns the private key of the current account
func (a *deviceKeystore) getAccountPrivateKey() (crypto.PrivKey, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.getOrGenerateNamedKey(keyAccount)
}

// getAccountProofPrivateKey returns the private proof key of
// the current account
func (a *deviceKeystore) getAccountProofPrivateKey() (crypto.PrivKey, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.getOrGenerateNamedKey(keyAccountProof)
}

// devicePrivateKey returns the current private key of the current device for
// the account and one-to-one conversations
func (a *deviceKeystore) devicePrivateKey() (crypto.PrivKey, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.getOrGenerateNamedKey(keyDevice)
}

// contactGroupPrivateKey retrieves the key for the contact group
// shared with the supplied contact's public key, this key will be derived to
// form the contact group keys
func (a *deviceKeystore) contactGroupPrivateKey(contactPublicKey crypto.PubKey) (crypto.PrivKey, error) {
	accountPrivateKey, err := a.getAccountPrivateKey()
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	return a.getOrComputeECDH(keyContactGroup, contactPublicKey, accountPrivateKey)
}

// memberDeviceForMultiMemberGroup retrieves the device private key for the
// supplied group
func (a *deviceKeystore) memberDeviceForMultiMemberGroup(groupPublicKey crypto.PubKey) (*ownMemberDevice, error) {
	memberPrivateKey, err := a.computeMemberKeyForMultiMemberGroup(groupPublicKey)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(fmt.Errorf("unable to get or generate a device key for group member: %w", err))
	}

	devicePrivateKey, err := a.getOrGenerateDeviceKeyForMultiMemberGroup(groupPublicKey)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	return newOwnMemberDevice(memberPrivateKey, devicePrivateKey), nil
}

// memberDeviceForGroup computes or retrieves the member and device key for the
// supplied group
func (a *deviceKeystore) memberDeviceForGroup(group *protocoltypes.Group) (*ownMemberDevice, error) {
	publicKey, err := group.GetPubKey()
	if err != nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("unable to get public key for group: %w", err))
	}

	switch group.GetGroupType() {
	case protocoltypes.GroupType_GroupTypeAccount, protocoltypes.GroupType_GroupTypeContact:
		memberPrivateKey, err := a.getAccountPrivateKey()
		if err != nil {
			return nil, errcode.ErrCode_ErrInternal.Wrap(err)
		}

		devicePrivateKey, err := a.devicePrivateKey()
		if err != nil {
			return nil, errcode.ErrCode_ErrInternal.Wrap(err)
		}

		return newOwnMemberDevice(memberPrivateKey, devicePrivateKey), nil

	case protocoltypes.GroupType_GroupTypeMultiMember:
		return a.memberDeviceForMultiMemberGroup(publicKey)
	}

	return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("unknown group type"))
}

// getOrGenerateNamedKey retrieves a private key by its name, or generate it
// if missing
func (a *deviceKeystore) getOrGenerateNamedKey(name string) (crypto.PrivKey, error) {
	privateKey, err := a.keystore.Get(name)
	if err == nil {
		return privateKey, nil
	} else if err.Error() != keystore.ErrNoSuchKey.Error() {
		return nil, errcode.ErrCode_ErrDBRead.Wrap(fmt.Errorf("unable to perform get operation on keystore: %w", err))
	}

	privateKey, _, err = crypto.GenerateEd25519Key(crand.Reader)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(fmt.Errorf("unable to generate an ed25519 key: %w", err))
	}

	if err := a.keystore.Put(name, privateKey); err != nil {
		return nil, errcode.ErrCode_ErrDBWrite.Wrap(fmt.Errorf("unable to perform put operation on keystore: %w", err))
	}

	return privateKey, nil
}

// getOrGenerateDeviceKeyForMultiMemberGroup fetches or generate a new device
// key for a multi-member group. The results do not need to be deterministic
// as it will only be used on the current device.
func (a *deviceKeystore) getOrGenerateDeviceKeyForMultiMemberGroup(groupPublicKey crypto.PubKey) (crypto.PrivKey, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	groupPublicKeyRaw, err := groupPublicKey.Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	name := strings.Join([]string{keyMemberDevice, hex.EncodeToString(groupPublicKeyRaw)}, "_")

	return a.getOrGenerateNamedKey(name)
}

// getOrComputeECDH fetches a named private key or computes one via an
// elliptic-curve Diffie-Hellman key agreement if not cached
func (a *deviceKeystore) getOrComputeECDH(nameSpace string, publicKey crypto.PubKey, ownPrivateKey crypto.PrivKey) (crypto.PrivKey, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	publicKeyRaw, err := publicKey.Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	name := strings.Join([]string{nameSpace, hex.EncodeToString(publicKeyRaw)}, "_")

	privateKey, err := a.keystore.Get(name)
	if err == nil {
		return privateKey, nil
	} else if err.Error() != keystore.ErrNoSuchKey.Error() {
		return nil, errcode.ErrCode_ErrDBRead.Wrap(fmt.Errorf("unable to perform get operation on keystore: %w", err))
	}

	privateKeyBytes, publicKeyBytes, err := cryptoutil.EdwardsToMontgomery(ownPrivateKey, publicKey)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoKeyConversion.Wrap(err)
	}

	secret := ecdh.X25519().ComputeSecret(privateKeyBytes, publicKeyBytes)
	groupSecretPrivateKey := ed25519.NewKeyFromSeed(secret)

	privateKey, _, err = crypto.KeyPairFromStdKey(&groupSecretPrivateKey)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoKeyConversion.Wrap(err)
	}

	if err := a.keystore.Put(name, privateKey); err != nil {
		return nil, errcode.ErrCode_ErrDBWrite.Wrap(err)
	}

	return privateKey, nil
}

// computeMemberKeyForMultiMemberGroup returns a deterministic private key
// for a multi member group, this allows a group to be joined from two
// different devices simultaneously without requiring a consensus.
func (a *deviceKeystore) computeMemberKeyForMultiMemberGroup(groupPublicKey crypto.PubKey) (crypto.PrivKey, error) {
	accountProofPrivateKey, err := a.getAccountProofPrivateKey()
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	return a.getOrComputeECDH(keyMember, groupPublicKey, accountProofPrivateKey)
}

// restoreAccountKeys restores exported LibP2P keys into the deviceKeystore, it
// will fail if accounts keys are already created or imported into the keystore
func (a *deviceKeystore) restoreAccountKeys(accountPrivateKeyBytes []byte, accountProofPrivateKeyBytes []byte) error {
	privateKeys := map[string]crypto.PrivKey{}

	for keyName, keyBytes := range map[string][]byte{
		keyAccount:      accountPrivateKeyBytes,
		keyAccountProof: accountProofPrivateKeyBytes,
	} {
		var err error
		privateKeys[keyName], err = getEd25519PrivateKeyFromLibP2PFormattedBytes(keyBytes)
		if err != nil {
			return errcode.ErrCode_ErrDeserialization.Wrap(err)
		}
	}

	if privateKeys[keyAccount].Equals(privateKeys[keyAccountProof]) {
		return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("the account key cannot be the same value as the account proof key"))
	}

	for keyName := range privateKeys {
		if exists, err := a.keystore.Has(keyName); err != nil {
			return errcode.ErrCode_ErrDBRead.Wrap(err)
		} else if exists {
			return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("an account is already set in this keystore"))
		}
	}

	for keyName, privateKey := range privateKeys {
		if err := a.keystore.Put(keyName, privateKey); err != nil {
			return errcode.ErrCode_ErrDBWrite.Wrap(err)
		}
	}

	return nil
}
