package secretstore

import (
	"context"

	"github.com/ipfs/go-cid"
	keystore "github.com/ipfs/go-ipfs-keystore"
	"github.com/libp2p/go-libp2p/core/crypto"
	"go.uber.org/zap"

	"berty.tech/weshnet/v2/pkg/protocoltypes"
)

const (
	PrecomputeOutOfStoreGroupRefsCount = 100
	PrecomputeMessageKeyCount          = 100
)

type messageKey [32]byte

type SecretStore interface {
	//
	// Account methods
	//

	// GetAccountProofPublicKey returns the user's account proof public key
	GetAccountProofPublicKey() (accountProofPublicKey crypto.PubKey, err error)

	// ImportAccountKeys restores backup of account keys into the SecretStore, it should fail if the store is already used by an account
	ImportAccountKeys(accountPrivateKey []byte, accountProofPrivateKey []byte) error

	// ExportAccountKeysForBackup returns the account's private key and proof private key of the user for a backup
	ExportAccountKeysForBackup() (accountPrivateKey []byte, accountProofPrivateKey []byte, err error)

	// GetAccountPrivateKey returns the account's private key, avoid using it, use GetGroupForAccount to get the account public key or sign data instead
	GetAccountPrivateKey() (accountPrivateKey crypto.PrivKey, err error)

	//
	// Groups methods
	//

	// GetGroupForAccount returns the Account's Group of the user
	GetGroupForAccount() (group *protocoltypes.Group, ownMemberDevice OwnMemberDevice, err error)

	// GetGroupForContact returns a contact group for communicating with the provided account
	GetGroupForContact(contactPublicKey crypto.PubKey) (group *protocoltypes.Group, err error)

	// PutGroup stores a group into the store
	PutGroup(ctx context.Context, group *protocoltypes.Group) error

	// FetchGroupByPublicKey gets an account from the store using the provided public key
	FetchGroupByPublicKey(ctx context.Context, publicKey crypto.PubKey) (group *protocoltypes.Group, err error)

	//
	// Envelopes methods
	//

	// OpenEnvelopeHeaders opens a message headers for a given group
	OpenEnvelopeHeaders(data []byte, group *protocoltypes.Group) (*protocoltypes.MessageEnvelope, *protocoltypes.MessageHeaders, error)

	// OpenEnvelopePayload opens a message payload with the given group headers
	OpenEnvelopePayload(ctx context.Context, msgEnvelope *protocoltypes.MessageEnvelope, msgHeaders *protocoltypes.MessageHeaders, groupPublicKey crypto.PubKey, ownPublicKey crypto.PubKey, msgCID cid.Cid) (*protocoltypes.EncryptedMessage, error)

	// SealEnvelope creates an encrypted payload to be sent to a group
	SealEnvelope(ctx context.Context, group *protocoltypes.Group, messagePayload []byte) (sealedEnvelope []byte, err error)

	//
	// Group member-device pairs methods
	//

	// GetOwnMemberDeviceForGroup gets a member and device key-pairs representing the current device in a given group
	GetOwnMemberDeviceForGroup(group *protocoltypes.Group) (OwnMemberDevice, error)

	//
	// Chain-keys methods
	//

	// RegisterChainKey records another device chain-key
	RegisterChainKey(ctx context.Context, group *protocoltypes.Group, senderDevicePublicKey crypto.PubKey, encryptedDeviceChainKey []byte) error

	// GetShareableChainKey returns a chain-key that can be decrypted by the provided member of a group
	GetShareableChainKey(ctx context.Context, group *protocoltypes.Group, targetMemberPublicKey crypto.PubKey) (encryptedDeviceChainKey []byte, err error)

	// IsChainKeyKnownForDevice checks whether a chain key of a device is already known
	IsChainKeyKnownForDevice(ctx context.Context, groupPublicKey crypto.PubKey, devicePublicKey crypto.PubKey) (isKnown bool)

	//
	// Out-of-store messages methods
	//

	// SealOutOfStoreMessageEnvelope encrypts a message to be sent outside a synchronized store
	SealOutOfStoreMessageEnvelope(id cid.Cid, env *protocoltypes.MessageEnvelope, headers *protocoltypes.MessageHeaders, group *protocoltypes.Group) (*protocoltypes.OutOfStoreMessageEnvelope, error)

	// OpenOutOfStoreMessage opens a message received outside a synchronized store
	OpenOutOfStoreMessage(ctx context.Context, payload []byte) (outOfStoreMessage *protocoltypes.OutOfStoreMessage, group *protocoltypes.Group, clearPayload []byte, alreadyDecrypted bool, err error)

	// UpdateOutOfStoreGroupReferences computes references of messages which might be received outside a synchronized store
	UpdateOutOfStoreGroupReferences(ctx context.Context, devicePublicKeyBytes []byte, first uint64, group *protocoltypes.Group) error

	// Close frees resources created by the secret store
	Close() error
}

// NewSecretStoreOptions contains the options that can be passed to NewSecretStore
type NewSecretStoreOptions struct {
	// PreComputedKeysCount specifies the number of keys to precompute,
	// defaults to PrecomputeMessageKeyCount
	PreComputedKeysCount int

	// PreComputedKeysCount specifies the number of out of store references
	// to precompute, defaults to PrecomputeOutOfStoreGroupRefsCount
	PrecomputeOutOfStoreGroupRefsCount int

	// Keystore specifies an implementation of a keystore to be used, can be
	// helpful if you want to rely on a hardware based keystore instead of a
	// software one
	Keystore keystore.Keystore

	// Logger specifies which logger to use, logging is disabled by default
	Logger *zap.Logger

	// DisableOutOfStoreSupport explicitly disables support of out-of-store
	// payloads
	DisableOutOfStoreSupport bool
}

// MemberDevice is the public keys of a device and its member
type MemberDevice interface {
	// Member returns the member public key
	Member() crypto.PubKey

	// Device returns the device public key
	Device() crypto.PubKey
}

// OwnMemberDevice is a MemberDevice for the current device, able to sign data
type OwnMemberDevice interface {
	MemberDevice

	// MemberSign signs the given data as a member of a group
	MemberSign(data []byte) ([]byte, error)

	// DeviceSign signs the given data as a device of a group
	DeviceSign(data []byte) ([]byte, error)
}
