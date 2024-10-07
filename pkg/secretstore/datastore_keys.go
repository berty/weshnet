package secretstore

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p/core/crypto"

	"berty.tech/weshnet/v2/pkg/errcode"
)

const (
	// dsNamespaceChainKeyForDeviceOnGroup is a namespace stores the current
	// state of a device chain key for a given group.
	// It contains the secret used to derive the next value of the chain key
	// and used to generate a message key for the message at `counter` value,
	// then put in the dsNamespacePrecomputedMessageKeys namespace.
	dsNamespaceChainKeyForDeviceOnGroup = "chainKeyForDeviceOnGroup"

	// dsNamespacePrecomputedMessageKeys is a namespace storing precomputed
	// message keys for a given group, device and message counter.
	// As the chain key stored has already been derived, these message keys
	// need to be computed beforehand.
	// The corresponding message can then be decrypted via a quick lookup.
	dsNamespacePrecomputedMessageKeys = "precomputedMessageKeys"

	// dsNamespaceMessageKeyForCIDs is a namespace containing the message key
	// for a given CID once the corresponding message has been decrypted.
	dsNamespaceMessageKeyForCIDs = "messageKeyForCIDs"

	// dsNamespaceOutOfStoreGroupHint is a namespace where HMAC value are
	// associated to a group public key.
	// It is used when receiving an out-of-store message (e.g. a push
	// notification) to identify the group on which the message belongs, which
	// can then be decrypted.
	dsNamespaceOutOfStoreGroupHint = "outOfStoreGroupHint"

	// dsNamespaceOutOfStoreGroupHintCounters is a namespace storing first and
	// last counter values for generated group hints inside the
	// dsNamespaceOutOfStoreGroupHint namespace
	dsNamespaceOutOfStoreGroupHintCounters = "outOfStoreGroupHintCounters"

	// dsNamespaceGroupDatastore is a namespace to store groups by their public
	// key
	dsNamespaceGroupDatastore = "groupByPublicKey"
)

func dsKeyForGroup(key []byte) datastore.Key {
	return datastore.KeyWithNamespaces([]string{
		dsNamespaceGroupDatastore,
		base64.RawURLEncoding.EncodeToString(key),
	})
}

// dsKeyForPrecomputedMessageKey returns a datastore.Key where will be stored a
// precalculated message key for a given group and device
func dsKeyForPrecomputedMessageKey(groupPublicKey, devicePublicKey []byte, counter uint64) datastore.Key {
	return datastore.KeyWithNamespaces([]string{
		dsNamespacePrecomputedMessageKeys,
		hex.EncodeToString(groupPublicKey),
		hex.EncodeToString(devicePublicKey),
		fmt.Sprintf("%d", counter),
	})
}

// dsKeyForCurrentChainKey returns a datastore.Key where will be stored a
// device chain key for a given group.
func dsKeyForCurrentChainKey(groupPublicKey crypto.PubKey, devicePublicKey crypto.PubKey) (datastore.Key, error) {
	devicePublicKeyBytes, err := devicePublicKey.Raw()
	if err != nil {
		return datastore.Key{}, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	groupPublicKeyBytes, err := groupPublicKey.Raw()
	if err != nil {
		return datastore.Key{}, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	return datastore.KeyWithNamespaces([]string{
		dsNamespaceChainKeyForDeviceOnGroup,
		hex.EncodeToString(groupPublicKeyBytes),
		hex.EncodeToString(devicePublicKeyBytes),
	}), nil
}

// dsKeyForMessageKeyByCID returns a datastore.Key where will be stored a
// message decryption key for a given message CID.
func dsKeyForMessageKeyByCID(id cid.Cid) datastore.Key {
	// TODO: specify the id
	return datastore.KeyWithNamespaces([]string{
		dsNamespaceMessageKeyForCIDs,
		id.String(),
	})
}

// dsKeyForOutOfStoreMessageGroupHint returns a datastore.Key where will be
// stored a group public key for a given push group reference.
func dsKeyForOutOfStoreMessageGroupHint(ref []byte) datastore.Key {
	return datastore.KeyWithNamespaces([]string{
		dsNamespaceOutOfStoreGroupHint,
		base64.RawURLEncoding.EncodeToString(ref),
	})
}

// dsKeyForOutOfStoreFirstLastCounters returns the datastore.Key where will be
// stored a protocoltypes.FirstLastCounters struct for the given group public
// key and device public key.
func dsKeyForOutOfStoreFirstLastCounters(groupPK, devicePK []byte) datastore.Key {
	return datastore.KeyWithNamespaces([]string{
		dsNamespaceOutOfStoreGroupHintCounters,
		base64.RawURLEncoding.EncodeToString(groupPK),
		base64.RawURLEncoding.EncodeToString(devicePK),
	})
}
