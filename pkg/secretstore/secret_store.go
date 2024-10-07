package secretstore

import (
	"context"
	"fmt"
	"sync"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	"github.com/libp2p/go-libp2p/core/crypto"
	"go.uber.org/zap"
	"golang.org/x/crypto/nacl/secretbox"
	"google.golang.org/protobuf/proto"

	"berty.tech/weshnet/v2/internal/datastoreutil"
	"berty.tech/weshnet/v2/pkg/cryptoutil"
	"berty.tech/weshnet/v2/pkg/errcode"
	"berty.tech/weshnet/v2/pkg/ipfsutil"
	"berty.tech/weshnet/v2/pkg/protocoltypes"
)

const (
	namespaceDeviceKeystore   = "device_keystore"
	namespaceOutOfStoreSecret = "push_secret_ref" // nolint:gosec
)

type secretStore struct {
	logger         *zap.Logger
	datastore      datastore.Datastore
	deviceKeystore *deviceKeystore

	messageMutex sync.RWMutex

	preComputedKeysCount               int
	precomputeOutOfStoreGroupRefsCount uint64
}

func (o *NewSecretStoreOptions) applyDefaults(rootDatastore datastore.Datastore) {
	if o.Logger == nil {
		o.Logger = zap.NewNop()
	}

	if o.Keystore == nil {
		o.Keystore = ipfsutil.NewDatastoreKeystore(datastoreutil.NewNamespacedDatastore(rootDatastore, datastore.NewKey(namespaceDeviceKeystore)))
	}

	if o.PreComputedKeysCount <= 0 {
		o.PreComputedKeysCount = PrecomputeMessageKeyCount
	}

	if o.PrecomputeOutOfStoreGroupRefsCount <= 0 {
		o.PrecomputeOutOfStoreGroupRefsCount = PrecomputeOutOfStoreGroupRefsCount
	}
}

// NewSecretStore instantiates a new SecretStore
func NewSecretStore(rootDatastore datastore.Datastore, opts *NewSecretStoreOptions) (SecretStore, error) {
	return newSecretStore(rootDatastore, opts)
}

// newSecretStore instantiates a new secretStore
func newSecretStore(rootDatastore datastore.Datastore, opts *NewSecretStoreOptions) (*secretStore, error) {
	if rootDatastore == nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("a datastore is required"))
	}

	if opts == nil {
		opts = &NewSecretStoreOptions{}
	}

	opts.applyDefaults(rootDatastore)

	devKeystore := newDeviceKeystore(opts.Keystore, opts.Logger)

	store := &secretStore{
		logger:         opts.Logger,
		datastore:      rootDatastore,
		deviceKeystore: devKeystore,

		preComputedKeysCount:               opts.PreComputedKeysCount,
		precomputeOutOfStoreGroupRefsCount: uint64(opts.PrecomputeOutOfStoreGroupRefsCount),
	}

	return store, nil
}

// NewInMemSecretStore instantiates a SecretStore using a volatile backend.
func NewInMemSecretStore(opts *NewSecretStoreOptions) (SecretStore, error) {
	return newInMemSecretStore(opts)
}

// newInMemSecretStore instantiates a secretStore using a volatile backend.
func newInMemSecretStore(opts *NewSecretStoreOptions) (*secretStore, error) {
	return newSecretStore(dssync.MutexWrap(datastore.NewMapDatastore()), opts)
}

func (s *secretStore) Close() error {
	return nil
}

func (s *secretStore) PutGroup(ctx context.Context, g *protocoltypes.Group) error {
	pk, err := g.GetPubKey()
	if err != nil {
		return errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	// TODO: check if partial group or full group and complete if necessary
	if ok, err := s.hasGroup(ctx, pk); err != nil {
		return errcode.ErrCode_ErrInvalidInput.Wrap(err)
	} else if ok {
		return nil
	}

	data, err := proto.Marshal(g)
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	if err := s.datastore.Put(ctx, dsKeyForGroup(g.GetPublicKey()), data); err != nil {
		return errcode.ErrCode_ErrKeystorePut.Wrap(err)
	}

	memberDevice, err := s.GetOwnMemberDeviceForGroup(g)
	if err != nil {
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	// Force generation of chain key for own device
	_, err = s.GetShareableChainKey(ctx, g, memberDevice.Member())
	if err != nil {
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	return nil
}

func (s *secretStore) GetOwnMemberDeviceForGroup(g *protocoltypes.Group) (OwnMemberDevice, error) {
	return s.deviceKeystore.memberDeviceForGroup(g)
}

func (s *secretStore) OpenOutOfStoreMessage(ctx context.Context, payload []byte) (*protocoltypes.OutOfStoreMessage, *protocoltypes.Group, []byte, bool, error) {
	oosMessageEnv := &protocoltypes.OutOfStoreMessageEnvelope{}
	if err := proto.Unmarshal(payload, oosMessageEnv); err != nil {
		return nil, nil, nil, false, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	groupPublicKey, err := s.OutOfStoreGetGroupPublicKeyByGroupReference(ctx, oosMessageEnv.GroupReference)
	if err != nil {
		return nil, nil, nil, false, errcode.ErrCode_ErrNotFound.Wrap(err)
	}

	oosMessage, err := s.decryptOutOfStoreMessageEnv(ctx, oosMessageEnv, groupPublicKey)
	if err != nil {
		return nil, nil, nil, false, errcode.ErrCode_ErrCryptoDecrypt.Wrap(err)
	}

	clear, newlyDecrypted, err := s.OutOfStoreMessageOpen(ctx, oosMessage, groupPublicKey)
	if err != nil {
		return nil, nil, nil, false, errcode.ErrCode_ErrCryptoDecrypt.Wrap(err)
	}

	group, err := s.FetchGroupByPublicKey(ctx, groupPublicKey)
	if err == nil {
		if err := s.UpdateOutOfStoreGroupReferences(ctx, oosMessage.DevicePk, oosMessage.Counter, group); err != nil {
			s.logger.Error("unable to update push group references", zap.Error(err))
		}
	}

	return oosMessage, group, clear, !newlyDecrypted, nil
}

func (s *secretStore) decryptOutOfStoreMessageEnv(ctx context.Context, env *protocoltypes.OutOfStoreMessageEnvelope, groupPK crypto.PubKey) (*protocoltypes.OutOfStoreMessage, error) {
	nonce, err := cryptoutil.NonceSliceToArray(env.Nonce)
	if err != nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	g, err := s.FetchGroupByPublicKey(ctx, groupPK)
	if err != nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("unable to find group, err: %w", err))
	}

	secret := g.GetSharedSecret()

	data, ok := secretbox.Open(nil, env.Box, nonce, secret)
	if !ok {
		return nil, errcode.ErrCode_ErrCryptoDecrypt.Wrap(fmt.Errorf("unable to decrypt message"))
	}

	outOfStoreMessage := &protocoltypes.OutOfStoreMessage{}
	if err := proto.Unmarshal(data, outOfStoreMessage); err != nil {
		return nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	return outOfStoreMessage, nil
}

func (s *secretStore) FetchGroupByPublicKey(ctx context.Context, publicKey crypto.PubKey) (*protocoltypes.Group, error) {
	keyBytes, err := publicKey.Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	data, err := s.datastore.Get(ctx, dsKeyForGroup(keyBytes))
	if err != nil {
		return nil, errcode.ErrCode_ErrMissingMapKey.Wrap(err)
	}

	g := &protocoltypes.Group{}
	if err := proto.Unmarshal(data, g); err != nil {
		return nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	return g, nil
}

func (s *secretStore) GetAccountProofPublicKey() (crypto.PubKey, error) {
	privateKey, err := s.deviceKeystore.getAccountPrivateKey()
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	return privateKey.GetPublic(), nil
}

func (s *secretStore) ImportAccountKeys(accountPrivateKeyBytes []byte, accountProofPrivateKeyBytes []byte) error {
	return s.deviceKeystore.restoreAccountKeys(accountPrivateKeyBytes, accountProofPrivateKeyBytes)
}

func (s *secretStore) ExportAccountKeysForBackup() (accountPrivateKeyBytes []byte, accountProofPrivateKeyBytes []byte, err error) {
	accountPrivateKey, err := s.deviceKeystore.getAccountPrivateKey()
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	accountProofPrivateKey, err := s.deviceKeystore.getAccountProofPrivateKey()
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	accountPrivateKeyBytes, err = crypto.MarshalPrivateKey(accountPrivateKey)
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	accountProofPrivateKeyBytes, err = crypto.MarshalPrivateKey(accountProofPrivateKey)
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	return accountPrivateKeyBytes, accountProofPrivateKeyBytes, nil
}

func (s *secretStore) GetAccountPrivateKey() (crypto.PrivKey, error) {
	accountPrivateKey, err := s.deviceKeystore.getAccountPrivateKey()
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	return accountPrivateKey, nil
}

func (s *secretStore) GetGroupForAccount() (*protocoltypes.Group, OwnMemberDevice, error) {
	accountPrivateKey, err := s.deviceKeystore.getAccountPrivateKey()
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrOrbitDBOpen.Wrap(err)
	}

	accountProofPrivateKey, err := s.deviceKeystore.getAccountProofPrivateKey()
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrOrbitDBOpen.Wrap(err)
	}

	devicePrivateKey, err := s.deviceKeystore.devicePrivateKey()
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	pubBytes, err := accountPrivateKey.GetPublic().Raw()
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	signingBytes, err := cryptoutil.SeedFromEd25519PrivateKey(accountProofPrivateKey)
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	return &protocoltypes.Group{
		PublicKey: pubBytes,
		Secret:    signingBytes,
		SecretSig: nil,
		GroupType: protocoltypes.GroupType_GroupTypeAccount,
	}, newOwnMemberDevice(accountPrivateKey, devicePrivateKey), nil
}

func (s *secretStore) GetGroupForContact(contactPublicKey crypto.PubKey) (*protocoltypes.Group, error) {
	contactPairPrivateKey, err := s.deviceKeystore.contactGroupPrivateKey(contactPublicKey)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	return getGroupForContact(contactPairPrivateKey)
}

func (s *secretStore) OpenEnvelopeHeaders(data []byte, g *protocoltypes.Group) (*protocoltypes.MessageEnvelope, *protocoltypes.MessageHeaders, error) {
	env := &protocoltypes.MessageEnvelope{}
	err := proto.Unmarshal(data, env)
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	nonce, err := cryptoutil.NonceSliceToArray(env.Nonce)
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrSerialization.Wrap(fmt.Errorf("unable to convert slice to array: %w", err))
	}

	headersBytes, ok := secretbox.Open(nil, env.MessageHeaders, nonce, g.GetSharedSecret())
	if !ok {
		return nil, nil, errcode.ErrCode_ErrCryptoDecrypt.Wrap(fmt.Errorf("secretbox failed to open headers"))
	}

	headers := &protocoltypes.MessageHeaders{}
	if err := proto.Unmarshal(headersBytes, headers); err != nil {
		return nil, nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	return env, headers, nil
}

func (s *secretStore) SealOutOfStoreMessageEnvelope(id cid.Cid, env *protocoltypes.MessageEnvelope, headers *protocoltypes.MessageHeaders, group *protocoltypes.Group) (*protocoltypes.OutOfStoreMessageEnvelope, error) {
	oosMessage := &protocoltypes.OutOfStoreMessage{
		Cid:              id.Bytes(),
		DevicePk:         headers.DevicePk,
		Counter:          headers.Counter,
		Sig:              headers.Sig,
		EncryptedPayload: env.Message,
		Nonce:            env.Nonce,
	}

	data, err := proto.Marshal(oosMessage)
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	nonce, err := cryptoutil.GenerateNonce()
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoNonceGeneration.Wrap(err)
	}

	secret, err := cryptoutil.KeySliceToArray(group.Secret)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoKeyConversion.Wrap(fmt.Errorf("unable to convert slice to array: %w", err))
	}

	encryptedData := secretbox.Seal(nil, data, nonce, secret)

	pushGroupRef, err := createOutOfStoreGroupReference(group, headers.DevicePk, headers.Counter)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	return &protocoltypes.OutOfStoreMessageEnvelope{
		Nonce:          nonce[:],
		Box:            encryptedData,
		GroupReference: pushGroupRef,
	}, nil
}

// hasGroup checks whether a group is already known by the secretStore
func (s *secretStore) hasGroup(ctx context.Context, key crypto.PubKey) (bool, error) {
	keyBytes, err := key.Raw()
	if err != nil {
		return false, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	return s.datastore.Has(ctx, dsKeyForGroup(keyBytes))
}
