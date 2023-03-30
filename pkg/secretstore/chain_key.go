package secretstore

import (
	crand "crypto/rand"
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"
	"golang.org/x/crypto/nacl/box"

	"berty.tech/weshnet/pkg/cryptoutil"
	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/protocoltypes"
)

// newDeviceChainKey creates a new random chain key
func newDeviceChainKey() (*protocoltypes.DeviceChainKey, error) {
	chainKey := make([]byte, 32)
	_, err := crand.Read(chainKey)
	if err != nil {
		return nil, errcode.ErrCryptoRandomGeneration.Wrap(err)
	}

	return &protocoltypes.DeviceChainKey{
		ChainKey: chainKey,
		Counter:  0,
	}, nil
}

// encryptDeviceChainKey encrypts a device chain key for a target member
func encryptDeviceChainKey(localDevicePrivateKey crypto.PrivKey, remoteMemberPubKey crypto.PubKey, deviceChainKey *protocoltypes.DeviceChainKey, group *protocoltypes.Group) ([]byte, error) {
	chainKeyBytes, err := deviceChainKey.Marshal()
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	mongPriv, mongPub, err := cryptoutil.EdwardsToMontgomery(localDevicePrivateKey, remoteMemberPubKey)
	if err != nil {
		return nil, errcode.ErrCryptoKeyConversion.Wrap(err)
	}

	nonce := groupIDToNonce(group)
	encryptedChainKey := box.Seal(nil, chainKeyBytes, nonce, mongPub, mongPriv)

	return encryptedChainKey, nil
}

// decryptDeviceChainKey decrypts a chain key sent by the given device
func decryptDeviceChainKey(encryptedDeviceChainKey []byte, group *protocoltypes.Group, localMemberPrivateKey crypto.PrivKey, senderDevicePubKey crypto.PubKey) (*protocoltypes.DeviceChainKey, error) {
	mongPriv, mongPub, err := cryptoutil.EdwardsToMontgomery(localMemberPrivateKey, senderDevicePubKey)
	if err != nil {
		return nil, errcode.ErrCryptoKeyConversion.Wrap(err)
	}

	nonce := groupIDToNonce(group)
	decryptedSecret := &protocoltypes.DeviceChainKey{}
	decryptedMessage, ok := box.Open(nil, encryptedDeviceChainKey, nonce, mongPub, mongPriv)
	if !ok {
		return nil, errcode.ErrCryptoDecrypt.Wrap(fmt.Errorf("unable to decrypt message"))
	}

	err = decryptedSecret.Unmarshal(decryptedMessage)
	if err != nil {
		return nil, errcode.ErrDeserialization.Wrap(err)
	}

	return decryptedSecret, nil
}

// groupIDToNonce converts a group public key to a value which can be used as
// a nonce of the nacl library
func groupIDToNonce(group *protocoltypes.Group) *[cryptoutil.NonceSize]byte {
	// Nonce doesn't need to be secret, random nor unpredictable, it just needs
	// to be used only once for a given sender+receiver set, and we will send
	// only one SecretEntryPayload per localDevicePrivateKey+remoteMemberPubKey
	// So we can reuse groupID as nonce for all SecretEntryPayload and save
	// 24 bytes of storage and bandwidth for each of them.
	//
	// See https://pynacl.readthedocs.io/en/stable/secret/#nonce
	// See Security Model here: https://nacl.cr.yp.to/box.html
	var nonce [cryptoutil.NonceSize]byte

	gid := group.GetPublicKey()

	copy(nonce[:], gid)

	return &nonce
}
