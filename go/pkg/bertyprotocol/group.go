package bertyprotocol

import (
	"context"
	"crypto/ed25519"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"

	"berty.tech/berty/v2/go/internal/cryptoutil"
	"berty.tech/berty/v2/go/pkg/bertytypes"
	"berty.tech/berty/v2/go/pkg/errcode"
	"berty.tech/go-orbit-db/events"
	"github.com/libp2p/go-libp2p-core/crypto"
	"go.uber.org/zap"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/nacl/box"
)

const CurrentGroupVersion = 1

// NewGroupMultiMember creates a new Group object and an invitation to be used by
// the first member of the group
func NewGroupMultiMember() (*bertytypes.Group, crypto.PrivKey, error) {
	priv, pub, err := crypto.GenerateEd25519Key(crand.Reader)
	if err != nil {
		return nil, nil, errcode.ErrCryptoKeyGeneration.Wrap(err)
	}

	pubBytes, err := pub.Raw()
	if err != nil {
		return nil, nil, errcode.ErrSerialization.Wrap(err)
	}

	signing, _, err := crypto.GenerateEd25519Key(crand.Reader)
	if err != nil {
		return nil, nil, errcode.ErrCryptoKeyGeneration.Wrap(err)
	}

	signingBytes, err := cryptoutil.SeedFromEd25519PrivateKey(signing)
	if err != nil {
		return nil, nil, errcode.ErrSerialization.Wrap(err)
	}

	skSig, err := priv.Sign(signingBytes)
	if err != nil {
		return nil, nil, errcode.ErrCryptoSignature.Wrap(err)
	}

	group := &bertytypes.Group{
		PublicKey: pubBytes,
		Secret:    signingBytes,
		SecretSig: skSig,
		GroupType: bertytypes.GroupTypeMultiMember,
	}

	return group, priv, nil
}

func getKeysForGroupOfContact(contactPairSK crypto.PrivKey) (crypto.PrivKey, crypto.PrivKey, error) {
	// Salt length must be equal to hash length (64 bytes for sha256)
	hash := sha256.New

	ck, err := contactPairSK.Raw()
	if err != nil {
		return nil, nil, errcode.ErrSerialization.Wrap(err)
	}

	// Generate Pseudo Random Key using ck as IKM and salt
	prk := hkdf.Extract(hash, ck, nil)
	if len(prk) == 0 {
		return nil, nil, errcode.ErrInternal
	}

	// Expand using extracted prk and groupID as info (kind of namespace)
	kdf := hkdf.Expand(hash, prk, nil)

	// Generate next KDF and message keys
	groupSeed, err := ioutil.ReadAll(io.LimitReader(kdf, 32))
	if err != nil {
		return nil, nil, errcode.ErrCryptoKeyGeneration.Wrap(err)
	}

	groupSecretSeed, err := ioutil.ReadAll(io.LimitReader(kdf, 32))
	if err != nil {
		return nil, nil, errcode.ErrCryptoKeyGeneration.Wrap(err)
	}

	sk1 := ed25519.NewKeyFromSeed(groupSeed)
	groupSK, _, err := crypto.KeyPairFromStdKey(&sk1)
	if err != nil {
		return nil, nil, errcode.ErrCryptoKeyGeneration.Wrap(err)
	}

	sk2 := ed25519.NewKeyFromSeed(groupSecretSeed)
	groupSecretSK, _, err := crypto.KeyPairFromStdKey(&sk2)
	if err != nil {
		return nil, nil, errcode.ErrCryptoKeyGeneration.Wrap(err)
	}

	return groupSK, groupSecretSK, nil
}

func getGroupForContact(contactPairSK crypto.PrivKey) (*bertytypes.Group, error) {
	groupSK, groupSecretSK, err := getKeysForGroupOfContact(contactPairSK)
	if err != nil {
		return nil, errcode.ErrCryptoKeyGeneration.Wrap(err)
	}
	pubBytes, err := groupSK.GetPublic().Raw()
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	signingBytes, err := cryptoutil.SeedFromEd25519PrivateKey(groupSecretSK)
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	return &bertytypes.Group{
		PublicKey: pubBytes,
		Secret:    signingBytes,
		SecretSig: nil,
		GroupType: bertytypes.GroupTypeContact,
	}, nil
}

func getGroupForAccount(priv, signing crypto.PrivKey) (*bertytypes.Group, error) {
	pubBytes, err := priv.GetPublic().Raw()
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	signingBytes, err := cryptoutil.SeedFromEd25519PrivateKey(signing)
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	return &bertytypes.Group{
		PublicKey: pubBytes,
		Secret:    signingBytes,
		SecretSig: nil,
		GroupType: bertytypes.GroupTypeAccount,
	}, nil
}

func metadataStoreListSecrets(ctx context.Context, gc *groupContext) map[crypto.PubKey]*bertytypes.DeviceSecret {
	publishedSecrets := map[crypto.PubKey]*bertytypes.DeviceSecret{}

	m := gc.MetadataStore()
	ownSK := gc.getMemberPrivKey()
	g := gc.Group()

	for meta := range m.ListEvents(ctx) {
		if meta == nil {
			continue
		}

		pk, ds, err := openDeviceSecret(meta.Metadata, ownSK, g)
		if errcode.Code(err) == errcode.ErrInvalidInput.Code() || errcode.Code(err) == errcode.ErrGroupSecretOtherDestMember.Code() {
			continue
		}

		if err != nil {
			gc.logger.Error("unable to open device secret", zap.Error(err))
			continue
		}

		publishedSecrets[pk] = ds
	}

	return publishedSecrets
}

func FillMessageKeysHolderUsingNewData(ctx context.Context, gc *groupContext) {
	m := gc.MetadataStore()

	for evt := range m.Subscribe(ctx) {
		go func(evt events.Event) {
			e, ok := evt.(*bertytypes.GroupMetadataEvent)
			if !ok {
				return
			}

			pk, ds, err := openDeviceSecret(e.Metadata, gc.getMemberPrivKey(), gc.Group())
			if errcode.Code(err) == errcode.ErrInvalidInput.Code() || errcode.Code(err) == errcode.ErrGroupSecretOtherDestMember.Code() {
				return
			}

			if err != nil {
				gc.logger.Error("an error occurred while opening device secrets", zap.Error(err))
				return
			}

			if err = gc.MessageKeystore().RegisterChainKey(gc.Group(), pk, ds, gc.DevicePubKey().Equals(pk)); err != nil {
				gc.logger.Error("unable to register chain key", zap.Error(err))
				return
			}
		}(evt)
	}
}

func FillMessageKeysHolderUsingPreviousData(ctx context.Context, gc *groupContext) {
	publishedSecrets := metadataStoreListSecrets(ctx, gc)
	wg := sync.WaitGroup{}

	for pk, sec := range publishedSecrets {
		wg.Add(1)
		go func(pk crypto.PubKey, sec *bertytypes.DeviceSecret) {
			if err := gc.MessageKeystore().RegisterChainKey(gc.Group(), pk, sec, gc.DevicePubKey().Equals(pk)); err != nil {
				gc.logger.Error("unable to register chain key", zap.Error(err))
			}
			wg.Done()
		}(pk, sec)
	}

	wg.Wait()
}

func ActivateGroupContext(ctx context.Context, gc *groupContext) error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	// Fill keystore
	go func() {
		wg.Done()
		FillMessageKeysHolderUsingNewData(ctx, gc)
	}()

	go func() {
		wg.Done()
		WatchNewMembersAndSendSecrets(ctx, gc.logger, gc)
	}()

	wg.Wait()
	wg.Add(2)

	go func() {
		start := time.Now()
		FillMessageKeysHolderUsingPreviousData(ctx, gc)

		gc.logger.Info(fmt.Sprintf("FillMessageKeysHolderUsingPreviousData took %s", time.Since(start)))
		wg.Done()
	}()

	go func() {
		start := time.Now()
		SendSecretsToExistingMembers(ctx, gc)

		gc.logger.Info(fmt.Sprintf("SendSecretsToExistingMembers took %s", time.Since(start)))
		wg.Done()
	}()

	wg.Wait()

	start := time.Now()
	if _, err := gc.MetadataStore().AddDeviceToGroup(ctx); err != nil {
		return errcode.ErrInternal.Wrap(err)
	}
	gc.logger.Info(fmt.Sprintf("AddDeviceToGroup took %s", time.Since(start)))

	return nil
}

func SendSecretsToExistingMembers(ctx context.Context, gctx *groupContext) {
	members := gctx.MetadataStore().ListMembers()
	wg := sync.WaitGroup{}

	for _, pk := range members {
		wg.Add(1)

		go func(pk crypto.PubKey) {
			defer wg.Done()

			rawPK, err := pk.Raw()
			if err != nil {
				gctx.logger.Error("failed to serialize pk", zap.Error(err))
			}

			if _, err := gctx.MetadataStore().SendSecret(ctx, pk); err != nil {
				if errcode.Code(err) != errcode.ErrGroupSecretAlreadySentToMember.Code() {
					gctx.logger.Info("secret already sent secret to member", zap.String("memberpk", base64.StdEncoding.EncodeToString(rawPK)))
					return
				}
			} else {
				gctx.logger.Info("sent secret to existing member", zap.String("memberpk", base64.StdEncoding.EncodeToString(rawPK)))
			}
		}(pk)
	}

	wg.Wait()
}

func WatchNewMembersAndSendSecrets(ctx context.Context, logger *zap.Logger, gctx *groupContext) {
	for evt := range gctx.MetadataStore().Subscribe(ctx) {
		go func(evt events.Event) {
			e, ok := evt.(*bertytypes.GroupMetadataEvent)
			if !ok {
				return
			}

			if e.Metadata.EventType != bertytypes.EventTypeGroupMemberDeviceAdded {
				return
			}

			event := &bertytypes.GroupAddMemberDevice{}
			if err := event.Unmarshal(e.Event); err != nil {
				logger.Error("unable to unmarshal payload", zap.Error(err))
				return
			}

			memberPK, err := crypto.UnmarshalEd25519PublicKey(event.MemberPK)
			if err != nil {
				logger.Error("unable to unmarshal sender member pk", zap.Error(err))
				return
			}

			if _, err := gctx.MetadataStore().SendSecret(ctx, memberPK); err != nil {
				if errcode.Code(err) != errcode.ErrGroupSecretAlreadySentToMember.Code() {
					logger.Error("unable to send secret to member", zap.Error(err))
				}
			}
		}(evt)
	}
}

func openDeviceSecret(m *bertytypes.GroupMetadata, localMemberPrivateKey crypto.PrivKey, group *bertytypes.Group) (crypto.PubKey, *bertytypes.DeviceSecret, error) {
	if m == nil || m.EventType != bertytypes.EventTypeGroupDeviceSecretAdded {
		return nil, nil, errcode.ErrInvalidInput
	}

	s := &bertytypes.GroupAddDeviceSecret{}
	if err := s.Unmarshal(m.Payload); err != nil {
		return nil, nil, errcode.ErrDeserialization.Wrap(err)
	}

	senderDevicePubKey, err := crypto.UnmarshalEd25519PublicKey(s.DevicePK)
	if err != nil {
		return nil, nil, errcode.ErrDeserialization.Wrap(err)
	}

	destMemberPubKey, err := crypto.UnmarshalEd25519PublicKey(s.DestMemberPK)
	if err != nil {
		return nil, nil, errcode.ErrDeserialization.Wrap(err)
	}

	if !localMemberPrivateKey.GetPublic().Equals(destMemberPubKey) {
		return nil, nil, errcode.ErrGroupSecretOtherDestMember
	}

	mongPriv, mongPub, err := cryptoutil.EdwardsToMontgomery(localMemberPrivateKey, senderDevicePubKey)
	if err != nil {
		return nil, nil, errcode.ErrCryptoKeyConversion.Wrap(err)
	}

	nonce := groupIDToNonce(group)
	decryptedSecret := &bertytypes.DeviceSecret{}
	decryptedMessage, ok := box.Open(nil, s.Payload, nonce, mongPub, mongPriv)
	if !ok {
		return nil, nil, errcode.ErrCryptoDecrypt
	}

	err = decryptedSecret.Unmarshal(decryptedMessage)
	if err != nil {
		return nil, nil, errcode.ErrDeserialization
	}

	return senderDevicePubKey, decryptedSecret, nil
}

func groupIDToNonce(group *bertytypes.Group) *[cryptoutil.NonceSize]byte {
	// Nonce doesn't need to be secret, random nor unpredictable, it just needs
	// to be used only once for a given {sender, receiver} set and we will send
	// only one SecretEntryPayload per {localDevicePrivKey, remoteMemberPubKey}
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
