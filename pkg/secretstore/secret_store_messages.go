package secretstore

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p/core/crypto"
	"go.uber.org/zap"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/nacl/secretbox"
	"google.golang.org/protobuf/proto"

	"berty.tech/weshnet/v2/pkg/cryptoutil"
	"berty.tech/weshnet/v2/pkg/errcode"
	"berty.tech/weshnet/v2/pkg/logutil"
	"berty.tech/weshnet/v2/pkg/protocoltypes"
)

// decryptionContext contains context about a decrypted message, its CID, the
// associated message key and whether it had been previously opened.
type decryptionContext struct {
	newlyDecrypted bool
	messageKey     *messageKey
	cid            cid.Cid
}

// computedMessageKey is a precomputed message key for a given counter used in the cache namespace.
type computedMessageKey struct {
	counter    uint64
	messageKey *messageKey
}

// getDeviceChainKeyForGroupAndDevice returns the device chain key for the given group and device.
func (s *secretStore) getDeviceChainKeyForGroupAndDevice(ctx context.Context, groupPublicKey crypto.PubKey, devicePublicKey crypto.PubKey) (*protocoltypes.DeviceChainKey, error) {
	if s == nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	key, err := dsKeyForCurrentChainKey(groupPublicKey, devicePublicKey)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	// Not mutex here
	dsBytes, err := s.datastore.Get(ctx, key)

	if err == datastore.ErrNotFound {
		return nil, errcode.ErrCode_ErrMissingInput.Wrap(err)
	} else if err != nil {
		return nil, errcode.ErrCode_ErrMessageKeyPersistenceGet.Wrap(err)
	}

	ds := &protocoltypes.DeviceChainKey{}
	if err := proto.Unmarshal(dsBytes, ds); err != nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	return ds, nil
}

// IsChainKeyKnownForDevice returns true if the device chain key is known for the given group and device.
func (s *secretStore) IsChainKeyKnownForDevice(ctx context.Context, groupPublicKey crypto.PubKey, devicePublicKey crypto.PubKey) (has bool) {
	if s == nil {
		return false
	}

	key, err := dsKeyForCurrentChainKey(groupPublicKey, devicePublicKey)
	if err != nil {
		return false
	}

	s.messageMutex.RLock()
	defer s.messageMutex.RUnlock()

	has, _ = s.datastore.Has(ctx, key)

	return
}

// delPrecomputedKey deletes the message key in the cache namespace for the given group, device and counter.
func (s *secretStore) delPrecomputedKey(ctx context.Context, groupPublicKey crypto.PubKey, devicePublicKey crypto.PubKey, msgCounter uint64) error {
	if s == nil {
		return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	devicePublicKeyRaw, err := devicePublicKey.Raw()
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	groupPublicKeyRaw, err := groupPublicKey.Raw()
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	id := dsKeyForPrecomputedMessageKey(groupPublicKeyRaw, devicePublicKeyRaw, msgCounter)

	err = s.datastore.Delete(ctx, id)
	if err != nil {
		return errcode.ErrCode_ErrMessageKeyPersistencePut.Wrap(err)
	}

	return nil
}

// postDecryptActions is called after a message has been decrypted.
// It saves the message key from the cache namespace to find it quickly on subsequent read operations.
// It derives the chain key in the cache namespace.
func (s *secretStore) postDecryptActions(ctx context.Context, decryptionCtx *decryptionContext, groupPublicKey crypto.PubKey, ownPublicKey crypto.PubKey, msgHeaders *protocoltypes.MessageHeaders) error {
	if s == nil {
		return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	// Message was newly decrypted, we can save the message key and derive
	// future keys if necessary.
	if decryptionCtx == nil || !decryptionCtx.newlyDecrypted {
		return nil
	}

	var (
		deviceChainKey *protocoltypes.DeviceChainKey
		err            error
	)

	devicePublicKey, err := crypto.UnmarshalEd25519PublicKey(msgHeaders.DevicePk)
	if err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	if err = s.putKeyForCID(ctx, decryptionCtx.cid, decryptionCtx.messageKey); err != nil {
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	if err = s.delPrecomputedKey(ctx, groupPublicKey, devicePublicKey, msgHeaders.Counter); err != nil {
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	if deviceChainKey, err = s.preComputeNextKey(ctx, groupPublicKey, devicePublicKey); err != nil {
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	// If the message was not emitted by the current Device we might need
	// to update the current chain key
	if ownPublicKey == nil || !ownPublicKey.Equals(devicePublicKey) {
		if err = s.updateCurrentKey(ctx, groupPublicKey, devicePublicKey, deviceChainKey); err != nil {
			return errcode.ErrCode_ErrInternal.Wrap(err)
		}
	}

	return nil
}

func (s *secretStore) GetShareableChainKey(ctx context.Context, group *protocoltypes.Group, targetMemberPublicKey crypto.PubKey) ([]byte, error) {
	deviceChainKey, err := s.getOwnDeviceChainKeyForGroup(ctx, group)
	if err != nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	privateMemberDevice, err := s.deviceKeystore.memberDeviceForGroup(group)
	if err != nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	encryptedDeviceChainKey, err := encryptDeviceChainKey(privateMemberDevice.device, targetMemberPublicKey, deviceChainKey, group)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoEncrypt.Wrap(err)
	}

	return encryptedDeviceChainKey, nil
}

// getOwnDeviceChainKeyForGroup returns the device chain key for the current
// device on a given group.
// If the chain key has not been created yet, it will be generated and
// registered.
func (s *secretStore) getOwnDeviceChainKeyForGroup(ctx context.Context, group *protocoltypes.Group) (*protocoltypes.DeviceChainKey, error) {
	if s == nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	if s.deviceKeystore == nil {
		return nil, errcode.ErrCode_ErrCryptoSignature.Wrap(fmt.Errorf("message keystore is opened in read-only mode"))
	}

	md, err := s.deviceKeystore.memberDeviceForGroup(group)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	groupPublicKey, err := group.GetPubKey()
	if err != nil {
		return nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	s.messageMutex.Lock()
	defer s.messageMutex.Unlock()

	ds, err := s.getDeviceChainKeyForGroupAndDevice(ctx, groupPublicKey, md.Device())
	if errcode.Is(err, errcode.ErrCode_ErrMissingInput) {
		// If secret does not exist, create it
		deviceChainKey, err := newDeviceChainKey()
		if err != nil {
			return nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
		}

		if err = s.registerChainKey(ctx, group, md.Device(), deviceChainKey, true); err != nil {
			return nil, errcode.ErrCode_ErrMessageKeyPersistencePut.Wrap(err)
		}

		return deviceChainKey, nil
	}
	if err != nil {
		return nil, errcode.ErrCode_ErrMessageKeyPersistenceGet.Wrap(err)
	}

	return ds, nil
}

// RegisterChainKey registers a chain key for the given group and device.
// If the device chain key is not from the current device, the function will
// precompute and store in the cache namespace the next message keys.
// It is the exported version of registerChainKey.
func (s *secretStore) RegisterChainKey(ctx context.Context, group *protocoltypes.Group, senderDevicePublicKey crypto.PubKey, encryptedDeviceChainKey []byte) error {
	if s == nil {
		return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	if s.deviceKeystore == nil {
		return errcode.ErrCode_ErrCryptoSignature.Wrap(fmt.Errorf("message keystore is opened in read-only mode"))
	}

	localMemberDevice, err := s.deviceKeystore.memberDeviceForGroup(group)
	if err != nil {
		return errcode.ErrCode_ErrGroupMemberUnknownGroupID.Wrap(err)
	}

	deviceChainKey, err := decryptDeviceChainKey(encryptedDeviceChainKey, group, localMemberDevice.member, senderDevicePublicKey)
	if err != nil {
		return errcode.ErrCode_ErrCryptoDecrypt.Wrap(err)
	}

	hasSecretBeenSentByCurrentDevice := localMemberDevice.Member().Equals(senderDevicePublicKey)

	return s.registerChainKey(ctx, group, senderDevicePublicKey, deviceChainKey, hasSecretBeenSentByCurrentDevice)
}

// registerChainKey registers a chain key for the given group and device.
// If the chain key is not from the current device, the function will
// precompute and store in the cache namespace the next message keys.
func (s *secretStore) registerChainKey(ctx context.Context, group *protocoltypes.Group, devicePublicKey crypto.PubKey, deviceChainKey *protocoltypes.DeviceChainKey, isCurrentDeviceChainKey bool) error {
	if s == nil {
		return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	groupPublicKey, err := group.GetPubKey()
	if err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	if _, err := s.getDeviceChainKeyForGroupAndDevice(ctx, groupPublicKey, devicePublicKey); err == nil {
		// Device is already registered, ignore it
		s.logger.Debug("device already registered in group",
			logutil.PrivateBinary("devicePublicKey", logutil.CryptoKeyToBytes(devicePublicKey)),
			logutil.PrivateBinary("groupPublicKey", logutil.CryptoKeyToBytes(groupPublicKey)),
		)
		return nil
	}

	s.logger.Debug("registering chain key",
		logutil.PrivateBinary("devicePublicKey", logutil.CryptoKeyToBytes(devicePublicKey)),
		logutil.PrivateBinary("groupPublicKey", logutil.CryptoKeyToBytes(groupPublicKey)),
	)

	// If own Device store key as is, no need to precompute future keys
	if isCurrentDeviceChainKey {
		if err := s.putDeviceChainKey(ctx, groupPublicKey, devicePublicKey, deviceChainKey); err != nil {
			return errcode.ErrCode_ErrInternal.Wrap(err)
		}

		return nil
	}

	s.messageMutex.Lock()

	if deviceChainKey, err = s.preComputeKeys(ctx, devicePublicKey, groupPublicKey, deviceChainKey); err != nil {
		s.messageMutex.Unlock()
		return errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	if err := s.putDeviceChainKey(ctx, groupPublicKey, devicePublicKey, deviceChainKey); err != nil {
		s.messageMutex.Unlock()
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	s.messageMutex.Unlock()

	devicePublicKeyBytes, err := devicePublicKey.Raw()
	if err == nil {
		if err := s.UpdateOutOfStoreGroupReferences(ctx, devicePublicKeyBytes, deviceChainKey.Counter, group); err != nil {
			s.logger.Error("updating out of store group references failed", zap.Error(err))
		}
	}

	return nil
}

// preComputeKeys precomputes the next m.preComputedKeysCount keys for the given device and group and put them in the cache namespace.
func (s *secretStore) preComputeKeys(ctx context.Context, devicePublicKey crypto.PubKey, groupPublicKey crypto.PubKey, deviceChainKey *protocoltypes.DeviceChainKey) (*protocoltypes.DeviceChainKey, error) {
	if s == nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	chainKeyValue := deviceChainKey.ChainKey
	counter := deviceChainKey.Counter

	groupPublicKeyBytes, err := groupPublicKey.Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	knownDeviceChainKey, err := s.getDeviceChainKeyForGroupAndDevice(ctx, groupPublicKey, devicePublicKey)
	if err != nil && !errcode.Is(err, errcode.ErrCode_ErrMissingInput) {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	var preComputedKeys []computedMessageKey
	for i := 0; i < s.getPrecomputedKeyExpectedCount(); i++ {
		counter++

		knownMK, err := s.getPrecomputedMessageKey(ctx, groupPublicKey, devicePublicKey, counter)
		if err != nil && !errcode.Is(err, errcode.ErrCode_ErrMissingInput) {
			return nil, errcode.ErrCode_ErrInternal.Wrap(err)
		}

		newChainKeyValue, mk, err := deriveNextKeys(chainKeyValue, nil, groupPublicKeyBytes)
		if err != nil {
			return nil, errcode.ErrCode_TODO.Wrap(err)
		}

		chainKeyValue = newChainKeyValue

		if knownMK != nil && knownDeviceChainKey != nil {
			if knownDeviceChainKey.Counter != counter-1 {
				continue
			}
		}

		preComputedKeys = append(preComputedKeys, computedMessageKey{counter, &mk})
	}

	err = s.putPrecomputedKeys(ctx, groupPublicKey, devicePublicKey, preComputedKeys)
	if err != nil {
		return nil, errcode.ErrCode_TODO.Wrap(err)
	}

	return &protocoltypes.DeviceChainKey{
		Counter:  counter,
		ChainKey: chainKeyValue,
	}, nil
}

// getPrecomputedKeyExpectedCount returns the number of precomputed keys that
// should be in the cache namespace of the keystore.
func (s *secretStore) getPrecomputedKeyExpectedCount() int {
	if s == nil || s.preComputedKeysCount < 0 {
		return 0
	}

	return s.preComputedKeysCount
}

// preComputeNextKey precomputes the next key for the given group and device and adds it to the cache namespace.
func (s *secretStore) preComputeNextKey(ctx context.Context, groupPublicKey crypto.PubKey, devicePublicKey crypto.PubKey) (*protocoltypes.DeviceChainKey, error) {
	if s == nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	if devicePublicKey == nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("devicePublicKey cannot be nil"))
	}

	groupPublicKeyBytes, err := groupPublicKey.Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	ds, err := s.getDeviceChainKeyForGroupAndDevice(ctx, groupPublicKey, devicePublicKey)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	newCounter := ds.Counter + 1

	// TODO: Salt?
	newCK, mk, err := deriveNextKeys(ds.ChainKey, nil, groupPublicKeyBytes)
	if err != nil {
		return nil, errcode.ErrCode_TODO.Wrap(err)
	}

	err = s.putPrecomputedKeys(ctx, groupPublicKey, devicePublicKey, []computedMessageKey{{newCounter, &mk}})
	if err != nil {
		return nil, errcode.ErrCode_TODO.Wrap(err)
	}

	return &protocoltypes.DeviceChainKey{
		Counter:  newCounter,
		ChainKey: newCK,
	}, nil
}

// getPrecomputedMessageKey returns the precomputed message key put in the cache
// namespace for the given group and device at the given counter.
func (s *secretStore) getPrecomputedMessageKey(ctx context.Context, groupPublicKey crypto.PubKey, devicePublicKey crypto.PubKey, counter uint64) (*messageKey, error) {
	if s == nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	deviceRaw, err := devicePublicKey.Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	groupRaw, err := groupPublicKey.Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	id := dsKeyForPrecomputedMessageKey(groupRaw, deviceRaw, counter)

	key, err := s.datastore.Get(ctx, id)

	if err == datastore.ErrNotFound {
		return nil, errcode.ErrCode_ErrMissingInput.Wrap(fmt.Errorf("key for message does not exist in datastore"))
	}
	if err != nil {
		return nil, errcode.ErrCode_ErrMessageKeyPersistenceGet.Wrap(err)
	}

	keyArray, err := cryptoutil.KeySliceToArray(key)
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	return (*messageKey)(keyArray), nil
}

// putPrecomputedKeys puts the given precomputed keys in the cache namespace.
// It will try to use a batch if the store supports it.
func (s *secretStore) putPrecomputedKeys(ctx context.Context, groupPublicKey crypto.PubKey, devicePublicKey crypto.PubKey, preComputedMessageKeys []computedMessageKey) error {
	if s == nil {
		return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	s.logger.Debug("putting precomputed keys", zap.Int("count", len(preComputedMessageKeys)))

	if len(preComputedMessageKeys) == 0 {
		return nil
	}

	deviceRaw, err := devicePublicKey.Raw()
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	groupRaw, err := groupPublicKey.Raw()
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	if batchedDatastore, ok := s.datastore.(datastore.BatchingFeature); ok {
		batch, err := batchedDatastore.Batch(ctx)
		if err == datastore.ErrBatchUnsupported {
			return s.putPrecomputedKeysNonBatched(ctx, groupRaw, deviceRaw, preComputedMessageKeys)
		}

		return s.putPrecomputedKeysBatched(ctx, batch, groupRaw, deviceRaw, preComputedMessageKeys)
	}

	return s.putPrecomputedKeysNonBatched(ctx, groupRaw, deviceRaw, preComputedMessageKeys)
}

func (s *secretStore) putPrecomputedKeysBatched(ctx context.Context, batch datastore.Batch, groupRaw []byte, deviceRaw []byte, preComputedMessageKeys []computedMessageKey) error {
	for _, preComputedKey := range preComputedMessageKeys {
		id := dsKeyForPrecomputedMessageKey(groupRaw, deviceRaw, preComputedKey.counter)

		if err := batch.Put(ctx, id, preComputedKey.messageKey[:]); err != nil {
			return errcode.ErrCode_ErrMessageKeyPersistencePut.Wrap(err)
		}
	}

	if err := batch.Commit(ctx); err != nil {
		return errcode.ErrCode_ErrMessageKeyPersistencePut.Wrap(err)
	}

	return nil
}

func (s *secretStore) putPrecomputedKeysNonBatched(ctx context.Context, groupRaw []byte, deviceRaw []byte, preComputedMessageKeys []computedMessageKey) error {
	for _, preComputedKey := range preComputedMessageKeys {
		id := dsKeyForPrecomputedMessageKey(groupRaw, deviceRaw, preComputedKey.counter)

		if err := s.datastore.Put(ctx, id, preComputedKey.messageKey[:]); err != nil {
			return errcode.ErrCode_ErrMessageKeyPersistencePut.Wrap(err)
		}
	}

	return nil
}

// putKeyForCID puts the given message key in the datastore for a specified CID.
func (s *secretStore) putKeyForCID(ctx context.Context, messageCID cid.Cid, messageKey *messageKey) error {
	if s == nil {
		return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	if !messageCID.Defined() {
		return nil
	}

	err := s.datastore.Put(ctx, dsKeyForMessageKeyByCID(messageCID), messageKey[:])
	if err != nil {
		return errcode.ErrCode_ErrMessageKeyPersistencePut.Wrap(err)
	}

	return nil
}

// OpenEnvelopePayload opens the payload of a message envelope and returns the
// decrypted message in its EncryptedMessage form.
// It also performs post decryption actions such as updating message key cache.
func (s *secretStore) OpenEnvelopePayload(ctx context.Context, msgEnvelope *protocoltypes.MessageEnvelope, msgHeaders *protocoltypes.MessageHeaders, groupPublicKey crypto.PubKey, ownDevicePublicKey crypto.PubKey, msgCID cid.Cid) (*protocoltypes.EncryptedMessage, error) {
	s.messageMutex.Lock()
	defer s.messageMutex.Unlock()

	msgBytes, decryptionCtx, err := s.openPayload(ctx, msgCID, groupPublicKey, msgEnvelope.Message, msgHeaders)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoDecryptPayload.Wrap(err)
	}

	if err := s.postDecryptActions(ctx, decryptionCtx, groupPublicKey, ownDevicePublicKey, msgHeaders); err != nil {
		return nil, errcode.ErrCode_TODO.Wrap(err)
	}

	var msg protocoltypes.EncryptedMessage
	err = proto.Unmarshal(msgBytes, &msg)
	if err != nil {
		return nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	return &msg, nil
}

// openPayload opens the payload of a message envelope and returns the
// decrypted message.
// It retrieves the message key from the keystore or the cache to decrypt
// the message.
func (s *secretStore) openPayload(ctx context.Context, msgCID cid.Cid, groupPublicKey crypto.PubKey, payload []byte, msgHeaders *protocoltypes.MessageHeaders) ([]byte, *decryptionContext, error) {
	if s == nil {
		return nil, nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	var (
		err           error
		decryptionCtx = &decryptionContext{
			cid:            msgCID,
			newlyDecrypted: true,
		}
		publicKey crypto.PubKey
	)

	if decryptionCtx.messageKey, err = s.getKeyForCID(ctx, msgCID); err == nil {
		decryptionCtx.newlyDecrypted = false
	} else {
		publicKey, err = crypto.UnmarshalEd25519PublicKey(msgHeaders.DevicePk)
		if err != nil {
			return nil, nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
		}

		decryptionCtx.messageKey, err = s.getPrecomputedMessageKey(ctx, groupPublicKey, publicKey, msgHeaders.Counter)
		if err != nil {
			return nil, nil, errcode.ErrCode_ErrCryptoDecrypt.Wrap(err)
		}
	}

	return s.openPayloadWithMessageKey(decryptionCtx, publicKey, payload, msgHeaders)
}

// openPayloadWithMessageKey opens the payload of a message envelope with the
// given key and returns the decrypted message with the decryptionContext
// struct.
func (s *secretStore) openPayloadWithMessageKey(decryptionCtx *decryptionContext, devicePublicKey crypto.PubKey, payload []byte, headers *protocoltypes.MessageHeaders) ([]byte, *decryptionContext, error) {
	msg, ok := secretbox.Open(nil, payload, uint64AsNonce(headers.Counter), (*[32]byte)(decryptionCtx.messageKey))
	if !ok {
		return nil, nil, errcode.ErrCode_ErrCryptoDecrypt.Wrap(fmt.Errorf("secret box failed to open message payload"))
	}

	if decryptionCtx.newlyDecrypted {
		if ok, err := devicePublicKey.Verify(msg, headers.Sig); !ok {
			return nil, nil, errcode.ErrCode_ErrCryptoSignatureVerification.Wrap(fmt.Errorf("unable to verify message signature"))
		} else if err != nil {
			return nil, nil, errcode.ErrCode_ErrCryptoSignatureVerification.Wrap(err)
		}
	}

	// Message was newly decrypted, we can save the message key and derive
	// future keys if necessary.
	return msg, decryptionCtx, nil
}

// getKeyForCID retrieves the message key for the given message CID.
func (s *secretStore) getKeyForCID(ctx context.Context, msgCID cid.Cid) (*messageKey, error) {
	if s == nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	if !msgCID.Defined() {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("undefined message CID"))
	}

	msgKey, err := s.datastore.Get(ctx, dsKeyForMessageKeyByCID(msgCID))

	if err == datastore.ErrNotFound {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	msgKeyArray, err := cryptoutil.KeySliceToArray(msgKey)
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	return (*messageKey)(msgKeyArray), nil
}

// putDeviceChainKey stores the chain key for the given group and device.
func (s *secretStore) putDeviceChainKey(ctx context.Context, groupPublicKey crypto.PubKey, devicePublicKey crypto.PubKey, deviceChainKey *protocoltypes.DeviceChainKey) error {
	if s == nil {
		return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	deviceChainKeyBytes, err := proto.Marshal(deviceChainKey)
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	datastoreKey, err := dsKeyForCurrentChainKey(groupPublicKey, devicePublicKey)
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	err = s.datastore.Put(ctx, datastoreKey, deviceChainKeyBytes)
	if err != nil {
		return errcode.ErrCode_ErrMessageKeyPersistencePut.Wrap(err)
	}

	return nil
}

// SealEnvelope encrypts the given payload and returns it as an envelope to be
// published on the group's store.
// It retrieves the device's chain key from the keystore to encrypt the payload
// using symmetric encryption. The payload is signed using the device's long
// term private key for the target group. It also updates the chain key and
// stores the next message key in the cache.
func (s *secretStore) SealEnvelope(ctx context.Context, group *protocoltypes.Group, messagePayload []byte) ([]byte, error) {
	if s == nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	if s.deviceKeystore == nil {
		return nil, errcode.ErrCode_ErrCryptoSignature.Wrap(fmt.Errorf("message keystore is opened in read-only mode"))
	}

	if group == nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("group cannot be nil"))
	}

	localMemberDevice, err := s.deviceKeystore.memberDeviceForGroup(group)
	if err != nil {
		return nil, errcode.ErrCode_ErrGroupMemberUnknownGroupID.Wrap(err)
	}

	groupPublicKey, err := group.GetPubKey()
	if err != nil {
		return nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	s.messageMutex.Lock()
	defer s.messageMutex.Unlock()

	deviceChainKey, err := s.getDeviceChainKeyForGroupAndDevice(ctx, groupPublicKey, localMemberDevice.Device())
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(fmt.Errorf("unable to get device chainkey: %w", err))
	}

	env, err := sealEnvelope(messagePayload, deviceChainKey, localMemberDevice.device, group)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoEncrypt.Wrap(fmt.Errorf("unable to seal envelope: %w", err))
	}

	if err := s.deriveDeviceChainKey(ctx, group, localMemberDevice.Device()); err != nil {
		return nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	return env, nil
}

// deriveDeviceChainKey derives the chain key value from the current one.
// It also updates the device chain key in the keystore.
func (s *secretStore) deriveDeviceChainKey(ctx context.Context, group *protocoltypes.Group, devicePublicKey crypto.PubKey) error {
	if s == nil {
		return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	if devicePublicKey == nil {
		return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("device public key cannot be nil"))
	}

	groupPublicKey, err := group.GetPubKey()
	if err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	deviceChainKey, err := s.preComputeNextKey(ctx, groupPublicKey, devicePublicKey)
	if err != nil {
		return errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	if err = s.updateCurrentKey(ctx, groupPublicKey, devicePublicKey, deviceChainKey); err != nil {
		return errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	return nil
}

// updateCurrentKey updates the current device chain key in the keystore if the
// given device secret has a higher counter.
func (s *secretStore) updateCurrentKey(ctx context.Context, groupPublicKey crypto.PubKey, devicePublicKey crypto.PubKey, deviceChainKey *protocoltypes.DeviceChainKey) error {
	if s == nil {
		return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	currentDeviceChainKey, err := s.getDeviceChainKeyForGroupAndDevice(ctx, groupPublicKey, devicePublicKey)
	if err != nil {
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	// FIXME: counter is set randomly and can overflow to 0
	if deviceChainKey.Counter < currentDeviceChainKey.Counter {
		return nil
	}

	if err = s.putDeviceChainKey(ctx, groupPublicKey, devicePublicKey, deviceChainKey); err != nil {
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	return nil
}

// OutOfStoreMessageOpen opens the given OutOfStoreMessage and returns the
// decrypted payload.
// The signature is verified against the given groupPublicKey.
// It derives the next message key and stores it in the cache, but it doesn't
// update the device's chain key.
func (s *secretStore) OutOfStoreMessageOpen(ctx context.Context, envelope *protocoltypes.OutOfStoreMessage, groupPublicKey crypto.PubKey) ([]byte, bool, error) {
	if s == nil {
		return nil, false, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("calling method of a non instantiated message keystore"))
	}

	if envelope == nil {
		return nil, false, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("envelope cannot be nil"))
	}

	if groupPublicKey == nil {
		return nil, false, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("group public key cannot be nil"))
	}

	devicePublicKey, err := crypto.UnmarshalEd25519PublicKey(envelope.DevicePk)
	if err != nil {
		return nil, false, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	c := cid.Undef
	if len(envelope.Cid) > 0 {
		_, c, err = cid.CidFromBytes(envelope.Cid)
		if err != nil {
			return nil, false, errcode.ErrCode_ErrDeserialization.Wrap(err)
		}
	}

	s.messageMutex.Lock()
	defer s.messageMutex.Unlock()

	decryptionCtx := &decryptionContext{newlyDecrypted: true}
	if decryptionCtx.messageKey, err = s.getKeyForCID(ctx, c); err == nil {
		decryptionCtx.newlyDecrypted = false
	} else {
		decryptionCtx.messageKey, err = s.getPrecomputedMessageKey(ctx, groupPublicKey, devicePublicKey, envelope.Counter)
		if err != nil {
			return nil, false, errcode.ErrCode_ErrCryptoDecrypt.Wrap(err)
		}
	}

	clear, decryptionCtx, err := s.openPayloadWithMessageKey(decryptionCtx, devicePublicKey, envelope.EncryptedPayload, &protocoltypes.MessageHeaders{
		Counter:  envelope.Counter,
		DevicePk: envelope.DevicePk,
		Sig:      envelope.Sig,
	})
	if err != nil {
		return nil, false, errcode.ErrCode_ErrCryptoDecrypt.Wrap(err)
	}

	if ok, err := devicePublicKey.Verify(clear, envelope.Sig); !ok {
		return nil, false, errcode.ErrCode_ErrCryptoSignatureVerification.Wrap(fmt.Errorf("unable to verify message signature"))
	} else if err != nil {
		return nil, false, errcode.ErrCode_ErrCryptoSignatureVerification.Wrap(err)
	}

	if _, err = s.preComputeNextKey(ctx, groupPublicKey, devicePublicKey); err != nil {
		return nil, false, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	return clear, decryptionCtx.newlyDecrypted, nil
}

// OutOfStoreGetGroupPublicKeyByGroupReference returns the group public key
// associated with the given out of store group reference (e.g. push
// notification payload).
func (s *secretStore) OutOfStoreGetGroupPublicKeyByGroupReference(ctx context.Context, ref []byte) (crypto.PubKey, error) {
	s.messageMutex.RLock()
	pk, err := s.datastore.Get(ctx, dsKeyForOutOfStoreMessageGroupHint(ref))
	s.messageMutex.RUnlock()

	if err != nil {
		return nil, errcode.ErrCode_ErrNotFound.Wrap(err)
	}

	groupPublicKey, err := crypto.UnmarshalEd25519PublicKey(pk)
	if err != nil {
		return nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	return groupPublicKey, nil
}

// UpdateOutOfStoreGroupReferences updates the out of store (e.g. push
// notification payload) group references for the given devicePublicKey and
// groupPublicKey in  the keystore. It creates the references for the
// given range [first + precomputeOutOfStoreGroupRefsCount] and
// [first - precomputeOutOfStoreGroupRefsCount] and deletes out of range
// references.
func (s *secretStore) UpdateOutOfStoreGroupReferences(ctx context.Context, devicePublicKey []byte, first uint64, group *protocoltypes.Group) error {
	s.messageMutex.Lock()
	defer s.messageMutex.Unlock()

	refsExisting := []uint64(nil)
	refsToCreate := []uint64(nil)

	currentFirst, currentLast, err := s.firstLastCachedGroupRefsForMember(ctx, devicePublicKey, group)
	if err == nil {
		for i := currentFirst; i != currentLast; i++ {
			refsExisting = append(refsExisting, i)
		}
	}

	// keep previous refs
	last := first + s.precomputeOutOfStoreGroupRefsCount
	first -= s.precomputeOutOfStoreGroupRefsCount
	for i := first; i != last; i++ {
		found := false

		// Ignore refs that should be kept
		for j := 0; j < len(refsExisting); j++ {
			if refsExisting[j] == i {
				refsExisting[j] = refsExisting[len(refsExisting)-1]
				refsExisting = refsExisting[:len(refsExisting)-1]
				found = true
				break
			}
		}

		if !found {
			refsToCreate = append(refsToCreate, i)
		}
	}

	// Remove useless old refs
	for i := 0; i < len(refsExisting); i++ {
		ref, err := createOutOfStoreGroupReference(group, devicePublicKey, refsExisting[i])
		if err != nil {
			s.logger.Error("creating existing out of store group reference failed", logutil.PrivateBinary("ref", ref), zap.Error(err))
			continue
		}

		if err := s.datastore.Delete(ctx, dsKeyForOutOfStoreMessageGroupHint(ref)); err != nil {
			s.logger.Error("deleting existing out of store group reference failed", logutil.PrivateBinary("ref", ref), zap.Error(err))
			continue
		}
	}

	// Add new refs
	for i := 0; i < len(refsToCreate); i++ {
		ref, err := createOutOfStoreGroupReference(group, devicePublicKey, refsToCreate[i])
		if err != nil {
			s.logger.Error("creating new out of store group reference failed", logutil.PrivateBinary("ref", ref), zap.Error(err))
			continue
		}

		if err := s.datastore.Put(ctx, dsKeyForOutOfStoreMessageGroupHint(ref), group.GetPublicKey()); err != nil {
			s.logger.Error("putting new out of store group reference failed", logutil.PrivateBinary("ref", ref), zap.Error(err))
			continue
		}
	}

	// Update first/last
	if err := s.putFirstLastCachedGroupRefsForMember(ctx, first, last, devicePublicKey, group); err != nil {
		s.logger.Error("putting first/last out of store group reference failed", zap.Error(err))
	}

	return nil
}

// firstLastCachedGroupRefsForMember returns the first and last cached group
// references counter for the given devicePublicKey and group.
func (s *secretStore) firstLastCachedGroupRefsForMember(ctx context.Context, devicePublicKeyBytes []byte, group *protocoltypes.Group) (uint64, uint64, error) {
	key := dsKeyForOutOfStoreFirstLastCounters(group.GetPublicKey(), devicePublicKeyBytes)

	// No mutex here
	bytes, err := s.datastore.Get(ctx, key)
	if err != nil {
		return 0, 0, errcode.ErrCode_ErrDBRead.Wrap(err)
	}

	ret := protocoltypes.FirstLastCounters{}
	if err := proto.Unmarshal(bytes, &ret); err != nil {
		return 0, 0, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	return ret.First, ret.Last, nil
}

// putFirstLastCachedGroupRefsForMember puts the first and last cached group
// references counter for the given devicePK and groupPK.
func (s *secretStore) putFirstLastCachedGroupRefsForMember(ctx context.Context, first uint64, last uint64, devicePublicKey []byte, group *protocoltypes.Group) error {
	key := dsKeyForOutOfStoreFirstLastCounters(group.GetPublicKey(), devicePublicKey)

	fistLast := protocoltypes.FirstLastCounters{
		First: first,
		Last:  last,
	}
	bytes, err := proto.Marshal(&fistLast)
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	// Not mutex here
	return s.datastore.Put(ctx, key, bytes)
}

func sealPayload(payload []byte, ds *protocoltypes.DeviceChainKey, devicePrivateKey crypto.PrivKey, g *protocoltypes.Group) ([]byte, []byte, error) {
	var (
		msgKey [32]byte
		err    error
	)

	sig, err := devicePrivateKey.Sign(payload)
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrCryptoSignature.Wrap(err)
	}

	if _, msgKey, err = deriveNextKeys(ds.ChainKey, nil, g.GetPublicKey()); err != nil {
		return nil, nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	return secretbox.Seal(nil, payload, uint64AsNonce(ds.Counter+1), &msgKey), sig, nil
}

func sealEnvelope(messagePayload []byte, deviceChainKey *protocoltypes.DeviceChainKey, devicePrivateKey crypto.PrivKey, g *protocoltypes.Group) ([]byte, error) {
	encryptedPayload, sig, err := sealPayload(messagePayload, deviceChainKey, devicePrivateKey, g)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoEncrypt.Wrap(err)
	}

	devicePublicKeyRaw, err := devicePrivateKey.GetPublic().Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	h := &protocoltypes.MessageHeaders{
		Counter:  deviceChainKey.Counter + 1,
		DevicePk: devicePublicKeyRaw,
		Sig:      sig,
	}

	headers, err := proto.Marshal(h)
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	nonce, err := cryptoutil.GenerateNonce()
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoNonceGeneration.Wrap(err)
	}

	encryptedHeaders := secretbox.Seal(nil, headers, nonce, g.GetSharedSecret())

	env, err := proto.Marshal(&protocoltypes.MessageEnvelope{
		MessageHeaders: encryptedHeaders,
		Message:        encryptedPayload,
		Nonce:          nonce[:],
	})
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	return env, nil
}

// nolint:unparam
func deriveNextKeys(chainKeyValue []byte, salt []byte, groupID []byte) ([]byte, messageKey, error) {
	var (
		nextMsg [32]byte
		err     error
	)

	// Salt length must be equal to hash length (64 bytes for sha256)
	hash := sha256.New

	// Generate Pseudo Random Key using chainKeyValue as IKM and salt
	prk := hkdf.Extract(hash, chainKeyValue, salt)
	if len(prk) == 0 {
		return nil, nextMsg, errcode.ErrCode_ErrInternal.Wrap(fmt.Errorf("unable to instantiate pseudo random key"))
	}

	// Expand using extracted prk and groupID as info (kind of namespace)
	kdf := hkdf.Expand(hash, prk, groupID)

	// Generate next KDF and message keys
	nextCK, err := io.ReadAll(io.LimitReader(kdf, 32))
	if err != nil {
		return nil, nextMsg, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	nextMsgSlice, err := io.ReadAll(io.LimitReader(kdf, 32))
	if err != nil {
		return nil, nextMsg, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	copy(nextMsg[:], nextMsgSlice)

	return nextCK, nextMsg, nil
}

func uint64AsNonce(val uint64) *[24]byte {
	var nonce [24]byte

	binary.BigEndian.PutUint64(nonce[:], val)

	return &nonce
}
