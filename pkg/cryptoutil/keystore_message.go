package cryptoutil

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	"github.com/libp2p/go-libp2p/core/crypto"
	"go.uber.org/zap"
	"golang.org/x/crypto/nacl/secretbox"

	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/logutil"
	"berty.tech/weshnet/pkg/protocoltypes"
)

const precomputePushRefsCount = 100

// MessageKeystore is a key-value store for storing values related to message
// opening. It has the following namespaces:
//   - `chainKeyForDeviceOnGroup`:
//     Storing the current state of a device chain key for a given group.
//     It contains the secret used to derive the next value of the chain key
//     and used to generate a message key for the message at `counter` value,
//     then put in the `precomputedMessageKeys` namespace.
//   - `precomputedMessageKeys`:
//     Storing precomputed message keys for a given group, device and message
//     counter. As the chain key stored has already been derived, these
//     message keys need to be computed beforehand. The corresponding message
//     can then be decrypted via a quick lookup.
//   - `messageKeyForCIDs`:
//     Containing the message key for a given message CID once the
//     corresponding message has been decrypted.
//   - `outOfStoreGroupHint`:
//     Keys are a HMAC value associated to a group public key. It is used when
//     receiving an out-of-store message (e.g. a push notification) to
//     identify the group on which the message belongs, which can then
//     be decrypted.
type MessageKeystore struct {
	lock                 sync.Mutex
	preComputedKeysCount int
	store                *dssync.MutexDatastore
	logger               *zap.Logger
}

type DecryptInfo struct {
	NewlyDecrypted bool
	MK             *[32]byte
	Cid            cid.Cid
}

// GetDeviceChainKey returns the device secret for the given group and device.
func (m *MessageKeystore) GetDeviceChainKey(ctx context.Context, groupPK, pk crypto.PubKey) (*protocoltypes.DeviceSecret, error) {
	if m == nil {
		return nil, errcode.ErrInvalidInput
	}

	pkB, err := pk.Raw()
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	groupRaw, err := groupPK.Raw()
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	key := idForCurrentCK(groupRaw, pkB)

	dsBytes, err := m.store.Get(ctx, key)
	if err == datastore.ErrNotFound {
		return nil, errcode.ErrMissingInput.Wrap(err)
	} else if err != nil {
		return nil, errcode.ErrMessageKeyPersistenceGet.Wrap(err)
	}

	ds := &protocoltypes.DeviceSecret{}
	if err := ds.Unmarshal(dsBytes); err != nil {
		return nil, errcode.ErrInvalidInput
	}

	return ds, nil
}

// HasSecretForRawDevicePK returns true if the device secret is known for the given group and device.
func (m *MessageKeystore) HasSecretForRawDevicePK(ctx context.Context, groupPK, devicePK []byte) (has bool) {
	if m == nil {
		return false
	}

	key := idForCurrentCK(groupPK, devicePK)
	has, _ = m.store.Has(ctx, key)
	return
}

// delPrecomputedKey deletes the message key in the cache namespace for the given group, device and counter.
func (m *MessageKeystore) delPrecomputedKey(ctx context.Context, groupPK, device crypto.PubKey, counter uint64) error {
	if m == nil {
		return errcode.ErrInvalidInput
	}

	deviceRaw, err := device.Raw()
	if err != nil {
		return errcode.ErrSerialization.Wrap(err)
	}

	groupRaw, err := groupPK.Raw()
	if err != nil {
		return errcode.ErrSerialization.Wrap(err)
	}

	id := idForPrecomputeMK(groupRaw, deviceRaw, counter)
	if err := m.store.Delete(ctx, id); err != nil {
		return errcode.ErrMessageKeyPersistencePut.Wrap(err)
	}

	return nil
}

// PostDecryptActions is called after a message has been decrypted.
// It saves the message key from the cache namespace to find it quickly on subsequent read operations.
// It derives the chain key in the cache namespace.
func (m *MessageKeystore) PostDecryptActions(ctx context.Context, di *DecryptInfo, g *protocoltypes.Group, ownPK crypto.PubKey, headers *protocoltypes.MessageHeaders) error {
	if m == nil {
		return errcode.ErrInvalidInput
	}

	// Message was newly decrypted, we can save the message key and derive
	// future keys if necessary.
	if di == nil || !di.NewlyDecrypted {
		return nil
	}

	var (
		ds  *protocoltypes.DeviceSecret
		err error
	)

	pk, err := crypto.UnmarshalEd25519PublicKey(headers.DevicePK)
	if err != nil {
		return errcode.ErrDeserialization.Wrap(err)
	}

	groupPK, err := g.GetPubKey()
	if err != nil {
		return errcode.ErrDeserialization.Wrap(err)
	}

	if err = m.putKeyForCID(ctx, di.Cid, di.MK); err != nil {
		return errcode.ErrInternal.Wrap(err)
	}

	if err = m.delPrecomputedKey(ctx, groupPK, pk, headers.Counter); err != nil {
		return errcode.ErrInternal.Wrap(err)
	}

	if ds, err = m.preComputeNextKey(ctx, groupPK, pk); err != nil {
		return errcode.ErrInternal.Wrap(err)
	}

	// If the message was not emitted by the current Device we might need
	// to update the current chain key
	if ownPK == nil || !ownPK.Equals(pk) {
		if err = m.updateCurrentKey(ctx, groupPK, pk, ds); err != nil {
			return errcode.ErrInternal.Wrap(err)
		}
	}

	return nil
}

// GetDeviceSecret returns the device secret for the current device on a given group.
// If the chain key has not been created yet, it will be generated and registered.
func (m *MessageKeystore) GetDeviceSecret(ctx context.Context, g *protocoltypes.Group, acc DeviceKeystore) (*protocoltypes.DeviceSecret, error) {
	if m == nil {
		return nil, errcode.ErrInvalidInput
	}

	md, err := acc.MemberDeviceForGroup(g)
	if err != nil {
		return nil, errcode.ErrInternal.Wrap(err)
	}

	groupPK, err := g.GetPubKey()
	if err != nil {
		return nil, errcode.ErrDeserialization.Wrap(err)
	}

	ds, err := m.GetDeviceChainKey(ctx, groupPK, md.device.GetPublic())
	if errcode.Is(err, errcode.ErrMissingInput) {
		// If secret does not exist, create it
		ds, err := NewDeviceSecret()
		if err != nil {
			return nil, errcode.ErrCryptoKeyGeneration.Wrap(err)
		}

		if err = m.registerChainKey(ctx, g, md.device.GetPublic(), ds, true); err != nil {
			return nil, errcode.ErrMessageKeyPersistencePut.Wrap(err)
		}

		return ds, nil
	}
	if err != nil {
		return nil, errcode.ErrMessageKeyPersistenceGet.Wrap(err)
	}

	return ds, nil
}

// RegisterChainKey registers a device secret for the given group and device.
// If the device secret is not from the current device, the function will
// precompute and store in the cache namespace the next message keys.
// It is the exported version of registerChainKey.
func (m *MessageKeystore) RegisterChainKey(ctx context.Context, g *protocoltypes.Group, devicePK crypto.PubKey, ds *protocoltypes.DeviceSecret, isOwnPK bool) error {
	if m == nil {
		return errcode.ErrInvalidInput
	}

	return m.registerChainKey(ctx, g, devicePK, ds, isOwnPK)
}

// registerChainKey registers a device secret for the given group and device.
// If the device secret is not from the current device, the function will
// precompute and store in the cache namespace the next message keys.
func (m *MessageKeystore) registerChainKey(ctx context.Context, g *protocoltypes.Group, devicePK crypto.PubKey, ds *protocoltypes.DeviceSecret, isOwnPK bool) error {
	if m == nil {
		return errcode.ErrInvalidInput
	}

	groupPK, err := g.GetPubKey()
	if err != nil {
		return errcode.ErrDeserialization.Wrap(err)
	}

	if _, err := m.GetDeviceChainKey(ctx, groupPK, devicePK); err == nil {
		// Device is already registered, ignore it
		m.logger.Debug("device already registered in group",
			logutil.PrivateBinary("devicePK", logutil.CryptoKeyToBytes(devicePK)),
			logutil.PrivateBinary("groupPK", logutil.CryptoKeyToBytes(groupPK)),
		)
		return nil
	}

	m.logger.Debug("registering chain key",
		logutil.PrivateBinary("devicePK", logutil.CryptoKeyToBytes(devicePK)),
		logutil.PrivateBinary("groupPK", logutil.CryptoKeyToBytes(groupPK)),
	)

	// If own Device store key as is, no need to precompute future keys
	if isOwnPK {
		if err := m.putDeviceChainKey(ctx, groupPK, devicePK, ds); err != nil {
			return errcode.ErrInternal.Wrap(err)
		}

		return nil
	}

	if ds, err = m.preComputeKeys(ctx, devicePK, groupPK, ds); err != nil {
		return errcode.ErrCryptoKeyGeneration.Wrap(err)
	}

	if err := m.putDeviceChainKey(ctx, groupPK, devicePK, ds); err != nil {
		return errcode.ErrInternal.Wrap(err)
	}

	devicePKBytes, err := devicePK.Raw()
	if err == nil {
		if err := m.UpdatePushGroupReferences(ctx, devicePKBytes, ds.Counter, g); err != nil {
			m.logger.Error("updating push group references failed", zap.Error(err))
		}
	}

	return nil
}

// preComputeKeys precomputes the next m.preComputedKeysCount keys for the given device and group and put them in the cache namespace.
func (m *MessageKeystore) preComputeKeys(ctx context.Context, device crypto.PubKey, groupPK crypto.PubKey, ds *protocoltypes.DeviceSecret) (*protocoltypes.DeviceSecret, error) {
	if m == nil {
		return nil, errcode.ErrInvalidInput
	}

	ck := ds.ChainKey
	counter := ds.Counter

	groupPKBytes, err := groupPK.Raw()
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	knownCK, err := m.GetDeviceChainKey(ctx, groupPK, device)
	if err != nil && !errcode.Is(err, errcode.ErrMissingInput) {
		return nil, errcode.ErrInternal.Wrap(err)
	}

	preComputedKeys := []computedKey{}
	for i := 0; i < m.GetPrecomputedKeyExpectedCount(); i++ {
		counter++

		knownMK, err := m.getPrecomputedKey(ctx, groupPK, device, counter)
		if err != nil && !errcode.Is(err, errcode.ErrMissingInput) {
			return nil, errcode.ErrInternal.Wrap(err)
		}

		// TODO: Salt?
		newCK, mk, err := deriveNextKeys(ck, nil, groupPKBytes)
		if err != nil {
			return nil, errcode.TODO.Wrap(err)
		}

		ck = newCK

		if knownMK != nil && knownCK != nil {
			if knownCK.Counter != counter-1 {
				continue
			}
		}

		preComputedKeys = append(preComputedKeys, computedKey{counter, &mk})
	}

	err = m.putPrecomputedKeys(ctx, groupPK, device, preComputedKeys...)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	return &protocoltypes.DeviceSecret{
		Counter:  counter,
		ChainKey: ck,
	}, nil
}

// preComputeNextKey precomputes the next key for the given group and device and adds it to the cache namespace.
func (m *MessageKeystore) preComputeNextKey(ctx context.Context, groupPK, devicePK crypto.PubKey) (*protocoltypes.DeviceSecret, error) {
	if m == nil || devicePK == nil {
		return nil, errcode.ErrInvalidInput
	}

	groupPKBytes, err := groupPK.Raw()
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	ds, err := m.GetDeviceChainKey(ctx, groupPK, devicePK)
	if err != nil {
		return nil, errcode.ErrInternal.Wrap(err)
	}

	newCounter := ds.Counter + 1

	// TODO: Salt?
	newCK, mk, err := deriveNextKeys(ds.ChainKey, nil, groupPKBytes)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	err = m.putPrecomputedKeys(ctx, groupPK, devicePK, computedKey{newCounter, &mk})
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	return &protocoltypes.DeviceSecret{
		Counter:  newCounter,
		ChainKey: newCK,
	}, nil
}

// getPrecomputedKey returns the precomputed key put in the cache namespace for the given group and device at the given counter.
func (m *MessageKeystore) getPrecomputedKey(ctx context.Context, groupPK, device crypto.PubKey, counter uint64) (*[32]byte, error) {
	if m == nil {
		return nil, errcode.ErrInvalidInput
	}

	deviceRaw, err := device.Raw()
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	groupRaw, err := groupPK.Raw()
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	id := idForPrecomputeMK(groupRaw, deviceRaw, counter)

	key, err := m.store.Get(ctx, id)

	if err == datastore.ErrNotFound {
		return nil, errcode.ErrMissingInput.Wrap(fmt.Errorf("key for message does not exist in datastore"))
	}
	if err != nil {
		return nil, errcode.ErrMessageKeyPersistenceGet.Wrap(err)
	}

	keyArray, err := KeySliceToArray(key)
	if err != nil {
		return nil, errcode.ErrSerialization
	}

	return keyArray, nil
}

// computedKey is a precomputed message key for a given counter used in the cache namespace.
type computedKey struct {
	counter uint64
	mk      *[32]byte
}

// putPrecomputedKeys puts the given precomputed keys in the cache namespace.
// It will try to use a batch if the store supports it.
func (m *MessageKeystore) putPrecomputedKeys(ctx context.Context, groupPK, device crypto.PubKey, preComputedKeys ...computedKey) error {
	if m == nil {
		return errcode.ErrInvalidInput
	}

	m.logger.Debug("putting precomputed keys", zap.Int("count", len(preComputedKeys)))

	if len(preComputedKeys) != 0 {
		deviceRaw, err := device.Raw()
		if err != nil {
			return errcode.ErrSerialization.Wrap(err)
		}

		groupRaw, err := groupPK.Raw()
		if err != nil {
			return errcode.ErrSerialization.Wrap(err)
		}

		batch, err := m.store.Batch(ctx)
		if err == datastore.ErrBatchUnsupported {
			for _, preComputedKey := range preComputedKeys {
				id := idForPrecomputeMK(groupRaw, deviceRaw, preComputedKey.counter)

				if err := m.store.Put(ctx, id, preComputedKey.mk[:]); err != nil {
					return errcode.ErrMessageKeyPersistencePut.Wrap(err)
				}
			}

			return nil
		} else if err != nil {
			return errcode.ErrMessageKeyPersistencePut.Wrap(err)
		}

		for _, preComputedKey := range preComputedKeys {
			id := idForPrecomputeMK(groupRaw, deviceRaw, preComputedKey.counter)

			if err := batch.Put(ctx, id, preComputedKey.mk[:]); err != nil {
				return errcode.ErrMessageKeyPersistencePut.Wrap(err)
			}
		}

		if err := batch.Commit(ctx); err != nil {
			return errcode.ErrMessageKeyPersistencePut.Wrap(err)
		}
	}

	return nil
}

// putKeyForCID puts the given message key in the datastore for the a specified CID.
func (m *MessageKeystore) putKeyForCID(ctx context.Context, id cid.Cid, key *[32]byte) error {
	if m == nil {
		return errcode.ErrInvalidInput
	}

	if !id.Defined() {
		return nil
	}

	err := m.store.Put(ctx, idForMK(id), key[:])
	if err != nil {
		return errcode.ErrMessageKeyPersistencePut.Wrap(err)
	}

	return nil
}

// OpenEnvelopePayload opens the payload of a message envelope and returns the decrypted message in its EncryptedMessage form.
// It also performs post decryption actions such as updating message key cache.
func (m *MessageKeystore) OpenEnvelopePayload(
	ctx context.Context,
	env *protocoltypes.MessageEnvelope,
	headers *protocoltypes.MessageHeaders,
	g *protocoltypes.Group,
	ownPK crypto.PubKey,
	id cid.Cid,
) (*protocoltypes.EncryptedMessage, error) {
	gPK, err := g.GetPubKey()
	if err != nil {
		return nil, errcode.ErrDeserialization.Wrap(err)
	}

	msgBytes, decryptInfo, err := m.OpenPayload(ctx, id, gPK, env.Message, headers)
	if err != nil {
		return nil, errcode.ErrCryptoDecryptPayload.Wrap(err)
	}

	if err := m.PostDecryptActions(ctx, decryptInfo, g, ownPK, headers); err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	var msg protocoltypes.EncryptedMessage
	err = msg.Unmarshal(msgBytes)
	if err != nil {
		return nil, errcode.ErrDeserialization.Wrap(err)
	}

	return &msg, nil
}

// OpenPayload opens the payload of a message envelope and returns the decrypted message.
// It retrieves the message key from the keystore or the cache to decrypt the message.
func (m *MessageKeystore) OpenPayload(ctx context.Context, id cid.Cid, groupPK crypto.PubKey, payload []byte, headers *protocoltypes.MessageHeaders) ([]byte, *DecryptInfo, error) {
	if m == nil {
		return nil, nil, errcode.ErrInvalidInput
	}

	var (
		err error
		di  = &DecryptInfo{
			Cid:            id,
			NewlyDecrypted: true,
		}
		pk crypto.PubKey
	)

	if di.MK, err = m.GetKeyForCID(ctx, id); err == nil {
		di.NewlyDecrypted = false
	} else {
		pk, err = crypto.UnmarshalEd25519PublicKey(headers.DevicePK)
		if err != nil {
			return nil, nil, errcode.ErrDeserialization.Wrap(err)
		}

		di.MK, err = m.getPrecomputedKey(ctx, groupPK, pk, headers.Counter)
		if err != nil {
			return nil, nil, errcode.ErrCryptoDecrypt.Wrap(err)
		}
	}

	return m.openPayload(di, pk, payload, headers)
}

// openPayload opens the payload of a message envelope with the given key and returns the decrypted message with the DecryptInfo struct.
func (m *MessageKeystore) openPayload(di *DecryptInfo, pk crypto.PubKey, payload []byte, headers *protocoltypes.MessageHeaders) ([]byte, *DecryptInfo, error) {
	msg, ok := secretbox.Open(nil, payload, uint64AsNonce(headers.Counter), di.MK)
	if !ok {
		return nil, nil, errcode.ErrCryptoDecrypt.Wrap(fmt.Errorf("secret box failed to open message payload"))
	}

	if di.NewlyDecrypted {
		if ok, err := pk.Verify(msg, headers.Sig); !ok {
			return nil, nil, errcode.ErrCryptoSignatureVerification.Wrap(fmt.Errorf("unable to verify message signature"))
		} else if err != nil {
			return nil, nil, errcode.ErrCryptoSignatureVerification.Wrap(err)
		}
	}

	// Message was newly decrypted, we can save the message key and derive
	// future keys if necessary.
	return msg, di, nil
}

// GetKeyForCID retrieves the message key for the given message CID.
func (m *MessageKeystore) GetKeyForCID(ctx context.Context, id cid.Cid) (*[32]byte, error) {
	if m == nil {
		return nil, errcode.ErrInvalidInput
	}

	if !id.Defined() {
		return nil, errcode.ErrInvalidInput
	}

	key, err := m.store.Get(ctx, idForMK(id))
	if err == datastore.ErrNotFound {
		return nil, errcode.ErrInvalidInput
	}

	keyArray, err := KeySliceToArray(key)
	if err != nil {
		return nil, errcode.ErrSerialization
	}

	return keyArray, nil
}

// GetPrecomputedKeyExpectedCount returns the number of precomputed keys that should be in the cache namespace of the keystore.
func (m *MessageKeystore) GetPrecomputedKeyExpectedCount() int {
	if m == nil {
		return 0
	}

	return m.preComputedKeysCount
}

// putDeviceChainKey stores the given device secret for the given groupPK and devicePK.
func (m *MessageKeystore) putDeviceChainKey(ctx context.Context, groupPK, device crypto.PubKey, ds *protocoltypes.DeviceSecret) error {
	if m == nil {
		return errcode.ErrInvalidInput
	}

	deviceRaw, err := device.Raw()
	if err != nil {
		return errcode.ErrSerialization.Wrap(err)
	}

	groupRaw, err := groupPK.Raw()
	if err != nil {
		return errcode.ErrSerialization.Wrap(err)
	}

	key := idForCurrentCK(groupRaw, deviceRaw)

	data, err := ds.Marshal()
	if err != nil {
		return errcode.ErrSerialization.Wrap(err)
	}

	err = m.store.Put(ctx, key, data)
	if err != nil {
		return errcode.ErrMessageKeyPersistencePut.Wrap(err)
	}

	return nil
}

// SealEnvelope encrypts the given payload and returns it as an envelope to be published on the group's store.
// It retrieves the device's chain key from the keystore to encrypt the payload using symmetric encryption.
// The payload is signed using the device's long term private key for the target group.
// It also updates the device secret and stores the next message key in the cache.
func (m *MessageKeystore) SealEnvelope(ctx context.Context, g *protocoltypes.Group, deviceSK crypto.PrivKey, payload []byte) ([]byte, error) {
	if m == nil {
		return nil, errcode.ErrInvalidInput
	}

	if deviceSK == nil || g == nil || m == nil {
		return nil, errcode.ErrInvalidInput
	}

	groupPK, err := g.GetPubKey()
	if err != nil {
		return nil, errcode.ErrDeserialization.Wrap(err)
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	ds, err := m.GetDeviceChainKey(ctx, groupPK, deviceSK.GetPublic())
	if err != nil {
		return nil, errcode.ErrInternal.Wrap(fmt.Errorf("unable to get device chainkey: %w", err))
	}

	env, err := SealEnvelope(payload, ds, deviceSK, g)
	if err != nil {
		return nil, errcode.ErrCryptoEncrypt.Wrap(fmt.Errorf("unable to seal envelope: %w", err))
	}

	if err := m.DeriveDeviceSecret(ctx, g, deviceSK.GetPublic()); err != nil {
		return nil, errcode.ErrCryptoKeyGeneration.Wrap(err)
	}

	return env, nil
}

// DeriveDeviceSecret derives the next device secret from the current device secret and stores it in the cache.
// It also updates the device secret in the keystore.
func (m *MessageKeystore) DeriveDeviceSecret(ctx context.Context, g *protocoltypes.Group, devicePK crypto.PubKey) error {
	if m == nil || devicePK == nil {
		return errcode.ErrInvalidInput
	}

	groupPK, err := g.GetPubKey()
	if err != nil {
		return errcode.ErrDeserialization.Wrap(err)
	}

	ds, err := m.preComputeNextKey(ctx, groupPK, devicePK)
	if err != nil {
		return errcode.ErrCryptoKeyGeneration.Wrap(err)
	}

	if err = m.updateCurrentKey(ctx, groupPK, devicePK, ds); err != nil {
		return errcode.ErrCryptoKeyGeneration.Wrap(err)
	}

	return nil
}

// updateCurrentKey updates the current device secret in the keystore if the given device secret has a higher counter.
func (m *MessageKeystore) updateCurrentKey(ctx context.Context, groupPK, pk crypto.PubKey, ds *protocoltypes.DeviceSecret) error {
	if m == nil {
		return errcode.ErrInvalidInput
	}

	currentCK, err := m.GetDeviceChainKey(ctx, groupPK, pk)
	if err != nil {
		return errcode.ErrInternal.Wrap(err)
	}

	if ds.Counter < currentCK.Counter {
		return nil
	}

	if err = m.putDeviceChainKey(ctx, groupPK, pk, ds); err != nil {
		return errcode.ErrInternal.Wrap(err)
	}

	return nil
}

// NewMessageKeystore instantiate a new MessageKeystore
func NewMessageKeystore(s datastore.Datastore, logger *zap.Logger) *MessageKeystore {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &MessageKeystore{
		preComputedKeysCount: 100,
		store:                dssync.MutexWrap(s),
		logger:               logger.Named("message-ks"),
	}
}

// nolint:deadcode,unused // NewInMemMessageKeystore instantiate a new MessageKeystore, useful for testing
func NewInMemMessageKeystore(logger *zap.Logger) (*MessageKeystore, func()) {
	ds := dssync.MutexWrap(datastore.NewMapDatastore())

	return NewMessageKeystore(ds, logger), func() { _ = ds.Close() }
}

// OpenOutOfStoreMessage opens the given OutOfStoreMessage and returns the decrypted payload.
// The signature is verified against the given devicePK.
// It derives the next message key and stores it in the cache, but it doesn't update the device secret.
func (m *MessageKeystore) OpenOutOfStoreMessage(ctx context.Context, envelope *protocoltypes.OutOfStoreMessage, groupPublicKey []byte) ([]byte, bool, error) {
	if m == nil || envelope == nil || len(groupPublicKey) == 0 {
		return nil, false, errcode.ErrInvalidInput
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	gPK, err := crypto.UnmarshalEd25519PublicKey(groupPublicKey)
	if err != nil {
		return nil, false, errcode.ErrDeserialization.Wrap(err)
	}

	dPK, err := crypto.UnmarshalEd25519PublicKey(envelope.DevicePK)
	if err != nil {
		return nil, false, errcode.ErrDeserialization.Wrap(err)
	}

	_, c, err := cid.CidFromBytes(envelope.CID)
	if err != nil {
		return nil, false, errcode.ErrDeserialization.Wrap(err)
	}

	di := &DecryptInfo{NewlyDecrypted: true}
	if di.MK, err = m.GetKeyForCID(ctx, c); err == nil {
		di.NewlyDecrypted = false
	} else {
		di.MK, err = m.getPrecomputedKey(ctx, gPK, dPK, envelope.Counter)
		if err != nil {
			return nil, false, errcode.ErrCryptoDecrypt.Wrap(err)
		}
	}

	clear, di, err := m.openPayload(di, dPK, envelope.EncryptedPayload, &protocoltypes.MessageHeaders{
		Counter:  envelope.Counter,
		DevicePK: envelope.DevicePK,
		Sig:      envelope.Sig,
	})
	if err != nil {
		return nil, false, errcode.ErrCryptoDecrypt.Wrap(err)
	}

	if ok, err := dPK.Verify(clear, envelope.Sig); !ok {
		return nil, false, errcode.ErrCryptoSignatureVerification.Wrap(fmt.Errorf("unable to verify message signature"))
	} else if err != nil {
		return nil, false, errcode.ErrCryptoSignatureVerification.Wrap(err)
	}

	if _, err = m.preComputeNextKey(ctx, gPK, dPK); err != nil {
		return nil, false, errcode.ErrInternal.Wrap(err)
	}

	return clear, di.NewlyDecrypted, nil
}

// refKey returns the datastore key of the groupPK for the given push group reference.
func (m *MessageKeystore) refKey(ref []byte) datastore.Key {
	return datastore.KeyWithNamespaces([]string{
		"outOfStoreGroupHint", base64.RawURLEncoding.EncodeToString(ref),
	})
}

// refFirstLastKey returns the datastore key of the FirstLastCounters struct for the given groupPK and devicePK.
func (m *MessageKeystore) refFirstLastKey(groupPK, devicePK []byte) datastore.Key {
	return datastore.KeyWithNamespaces([]string{
		"outOfStoreGroupHint",
		base64.RawURLEncoding.EncodeToString(groupPK),
		base64.RawURLEncoding.EncodeToString(devicePK),
	})
}

// GetByPushGroupReference returns the groupPK associated with the given push group reference.
func (m *MessageKeystore) GetByPushGroupReference(ctx context.Context, ref []byte) ([]byte, error) {
	return m.store.Get(ctx, m.refKey(ref))
}

// UpdatePushGroupReferences updates the push group references for the given devicePK and groupPK in the keystore.
// It creates the references for the given range [first + precomputePushRefsCount] and [first - precomputePushRefsCount] and deletes the references out of range.
func (m *MessageKeystore) UpdatePushGroupReferences(ctx context.Context, devicePK []byte, first uint64, group GroupWithSecret) error {
	refsExisting := []uint64(nil)
	refsToCreate := []uint64(nil)

	groupPushSecret, err := GetGroupPushSecret(group)
	if err != nil {
		return errcode.ErrCryptoKeyGeneration.Wrap(err)
	}

	currentFirst, currentLast, err := m.firstLastCachedGroupRefsForMember(ctx, devicePK, group)
	if err == nil {
		for i := currentFirst; i != currentLast; i++ {
			refsExisting = append(refsExisting, i)
		}
	}

	// keep previous refs
	last := first + precomputePushRefsCount
	first -= precomputePushRefsCount
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
		ref, err := CreatePushGroupReference(devicePK, refsExisting[i], groupPushSecret)
		if err != nil {
			m.logger.Error("creating existing push group reference failed", logutil.PrivateBinary("ref", ref), zap.Error(err))
			continue
		}

		if err := m.store.Delete(ctx, m.refKey(ref)); err != nil {
			m.logger.Error("deleting existing push group reference failed", logutil.PrivateBinary("ref", ref), zap.Error(err))
			continue
		}
	}

	// Add new refs
	for i := 0; i < len(refsToCreate); i++ {
		ref, err := CreatePushGroupReference(devicePK, refsToCreate[i], groupPushSecret)
		if err != nil {
			m.logger.Error("creating new push group reference failed", logutil.PrivateBinary("ref", ref), zap.Error(err))
			continue
		}

		if err := m.store.Put(ctx, m.refKey(ref), group.GetPublicKey()); err != nil {
			m.logger.Error("putting new push group reference failed", logutil.PrivateBinary("ref", ref), zap.Error(err))
			continue
		}
	}

	// Update first/last
	if err := m.putFirstLastCachedGroupRefsForMember(ctx, first, last, devicePK, group); err != nil {
		m.logger.Error("putting first/last push group reference failed", zap.Error(err))
	}

	return nil
}

// firstLastCachedGroupRefsForMember returns the first and last cached group references counter for the given devicePK and groupPK.
func (m *MessageKeystore) firstLastCachedGroupRefsForMember(ctx context.Context, devicePK []byte, group GroupWithSecret) (uint64, uint64, error) {
	key := m.refFirstLastKey(group.GetPublicKey(), devicePK)
	bytes, err := m.store.Get(ctx, key)
	if err != nil {
		return 0, 0, err
	}

	ret := protocoltypes.FirstLastCounters{}
	if err := ret.Unmarshal(bytes); err != nil {
		return 0, 0, err
	}

	return ret.First, ret.Last, nil
}

// putFirstLastCachedGroupRefsForMember puts the first and last cached group references counter for the given devicePK and groupPK.
func (m *MessageKeystore) putFirstLastCachedGroupRefsForMember(ctx context.Context, first uint64, last uint64, devicePK []byte, group GroupWithSecret) error {
	key := m.refFirstLastKey(group.GetPublicKey(), devicePK)

	fistLast := protocoltypes.FirstLastCounters{
		First: first,
		Last:  last,
	}
	bytes, err := fistLast.Marshal()
	if err != nil {
		return err
	}

	return m.store.Put(ctx, key, bytes)
}
