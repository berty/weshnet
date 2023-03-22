package secretstore

import (
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"berty.tech/weshnet/pkg/protocoltypes"
)

func Test_PushGroupReferences(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, _, err := protocoltypes.NewGroupMultiMember()
	require.NoError(t, err)

	otherSecretStore, err := newInMemSecretStore(nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = otherSecretStore.Close()
	})

	ownSecretStore, err := newInMemSecretStore(nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = ownSecretStore.Close()
	})

	otherMemberDevice, err := otherSecretStore.GetOwnMemberDeviceForGroup(g)
	require.NoError(t, err)

	otherDevicePK, err := otherMemberDevice.Device().Raw()
	require.NoError(t, err)

	deviceChainKey, err := newDeviceChainKey()
	require.NoError(t, err)

	// test with the deviceChainKey counter
	updateAndTestPushGroupReferences(ctx, ownSecretStore, otherDevicePK, deviceChainKey.Counter, g, t)

	// do the same test with a new device chain key counter,
	// so we can test if old references are deleted
	updateAndTestPushGroupReferences(ctx, ownSecretStore, otherDevicePK, deviceChainKey.Counter+10, g, t)
}

func updateAndTestPushGroupReferences(ctx context.Context, secretStore *secretStore, devicePK []byte, counter uint64, g *protocoltypes.Group, t *testing.T) {
	// update the push group references
	err := secretStore.UpdateOutOfStoreGroupReferences(ctx, devicePK, counter, g)
	require.NoError(t, err)

	// test that the push group references are updated
	// refs start counter - 100 to counter + 100
	start := counter - PrecomputeOutOfStoreGroupRefsCount
	end := counter + PrecomputeOutOfStoreGroupRefsCount

	for i := start; i < end; i++ {
		// compute the push group reference
		pushGroupRef, err := createOutOfStoreGroupReference(g, devicePK, i)
		require.NoError(t, err)

		_, err = secretStore.OutOfStoreGetGroupPublicKeyByGroupReference(ctx, pushGroupRef)
		require.NoError(t, err, fmt.Sprintf("started at %d, failed as %d", start, i))
	}

	// test boundary conditions

	// before the start counter
	{
		before := counter - PrecomputeOutOfStoreGroupRefsCount - 1
		pushGroupRef, err := createOutOfStoreGroupReference(g, devicePK, before)
		require.NoError(t, err)
		_, err = secretStore.OutOfStoreGetGroupPublicKeyByGroupReference(ctx, pushGroupRef)
		require.Error(t, err)
	}

	// after the end counter
	{
		end := counter + PrecomputeOutOfStoreGroupRefsCount + 1
		pushGroupRef, err := createOutOfStoreGroupReference(g, devicePK, end)
		require.NoError(t, err)
		_, err = secretStore.OutOfStoreGetGroupPublicKeyByGroupReference(ctx, pushGroupRef)
		require.Error(t, err)
	}
}

func Test_SealOutOfStoreMessageEnvelope_OpenOutOfStoreMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, _, err := protocoltypes.NewGroupMultiMember()
	require.NoError(t, err)

	accUnrelated, err := newInMemSecretStore(nil)
	require.NoError(t, err)
	require.NotNil(t, accUnrelated)
	t.Cleanup(func() {
		_ = accUnrelated.Close()
	})

	err = accUnrelated.PutGroup(ctx, g)
	require.NoError(t, err)

	acc1, err := newInMemSecretStore(nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = acc1.Close()
	})

	acc2, err := newInMemSecretStore(nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = acc1.Close()
	})

	memberDevice1ForGroup, err := acc1.GetOwnMemberDeviceForGroup(g)
	require.NoError(t, err)

	memberDevice2ForGroup, err := acc2.GetOwnMemberDeviceForGroup(g)
	require.NoError(t, err)

	deviceChainKey1For2, err := acc1.GetShareableChainKey(ctx, g, memberDevice2ForGroup.Member())
	require.NoError(t, err)

	testPayload := []byte("test payload")

	err = acc2.RegisterChainKey(ctx, g, memberDevice1ForGroup.Device(), deviceChainKey1For2)
	require.NoError(t, err)

	envEncrypted, err := acc1.SealEnvelope(ctx, g, testPayload)
	require.NoError(t, err)

	env, headers, err := ((*secretStore)(nil)).OpenEnvelopeHeaders(envEncrypted, g)
	require.NoError(t, err)

	outOfStoreEnv, err := (*secretStore)(nil).SealOutOfStoreMessageEnvelope(cid.Undef, env, headers, g)
	require.NoError(t, err)

	groupPublicKey, err := acc2.OutOfStoreGetGroupPublicKeyByGroupReference(ctx, outOfStoreEnv.GroupReference)
	require.NoError(t, err)

	outOfStoreMessage, err := accUnrelated.decryptOutOfStoreMessageEnv(ctx, outOfStoreEnv, groupPublicKey)
	require.NoError(t, err)

	payload, newlyDecrypted, err := acc2.OutOfStoreMessageOpen(ctx, outOfStoreMessage, groupPublicKey)
	require.NoError(t, err)
	require.True(t, newlyDecrypted)
	require.Equal(t, testPayload, payload)
}

func Test_OutOfStoreDeserialize(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	secretStore1, err := NewSecretStore(datastore.NewMapDatastore(), nil)
	require.NoError(t, err)

	t.Cleanup(func() { _ = secretStore1.Close() })

	secretStore2, err := NewSecretStore(datastore.NewMapDatastore(), nil)
	require.NoError(t, err)

	t.Cleanup(func() { _ = secretStore2.Close() })

	group, _, err := protocoltypes.NewGroupMultiMember()
	require.NoError(t, err)

	require.NoError(t, secretStore1.PutGroup(ctx, group))
	require.NoError(t, secretStore2.PutGroup(ctx, group))

	memberDevice1, err := secretStore1.GetOwnMemberDeviceForGroup(group)
	require.NoError(t, err)

	memberDevice2, err := secretStore2.GetOwnMemberDeviceForGroup(group)
	require.NoError(t, err)

	chainKey1For2, err := secretStore1.GetShareableChainKey(ctx, group, memberDevice2.Member())
	require.NoError(t, err)

	chainKey2For1, err := secretStore2.GetShareableChainKey(ctx, group, memberDevice1.Member())
	require.NoError(t, err)

	require.NoError(t, secretStore1.RegisterChainKey(ctx, group, memberDevice2.Device(), chainKey2For1))
	require.NoError(t, secretStore2.RegisterChainKey(ctx, group, memberDevice1.Device(), chainKey1For2))

	testPayload := []byte("test payload")

	env, err := secretStore1.SealEnvelope(ctx, group, testPayload)
	require.NoError(t, err)

	envHeaders, msgHeaders, err := secretStore1.OpenEnvelopeHeaders(env, group)
	require.NoError(t, err)

	dummyCID, err := cid.Parse("QmNR2n4zywCV61MeMLB6JwPueAPqheqpfiA4fLPMxouEmQ")
	require.NoError(t, err)

	device1Raw, err := memberDevice1.Device().Raw()
	require.NoError(t, err)

	outOfStoreMessageEnvelope, err := secretStore1.SealOutOfStoreMessageEnvelope(dummyCID, envHeaders, msgHeaders, group)
	require.NoError(t, err)

	outOfStoreMessageEnvelopeBytes, err := outOfStoreMessageEnvelope.Marshal()
	require.NoError(t, err)

	// Attempting to decrypt the message without a relay
	{
		openedOutOfStoreMessage, groupFound, clearPayload, alreadyDecrypted, err := secretStore2.OpenOutOfStoreMessage(ctx, outOfStoreMessageEnvelopeBytes)
		require.NoError(t, err)

		require.Equal(t, testPayload, clearPayload)
		require.Equal(t, device1Raw, openedOutOfStoreMessage.DevicePK)
		require.Equal(t, dummyCID.Bytes(), openedOutOfStoreMessage.CID)
		require.Equal(t, group.PublicKey, groupFound.PublicKey)
		require.False(t, alreadyDecrypted)
	}
}

func Test_EncryptMessageEnvelopeAndDerive(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, _, err := protocoltypes.NewGroupMultiMember()
	assert.NoError(t, err)

	mkh1, err := newInMemSecretStore(nil)
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = mkh1.Close()
	})

	mkh2, err := newInMemSecretStore(nil)
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = mkh1.Close()
	})

	omd1, err := mkh1.GetOwnMemberDeviceForGroup(g)
	assert.NoError(t, err)

	omd2, err := mkh2.GetOwnMemberDeviceForGroup(g)
	assert.NoError(t, err)

	gPK, err := g.GetPubKey()
	assert.NoError(t, err)

	gc1DevicePubKey := omd1.Device()

	ds1For1Encrypted, err := mkh1.GetShareableChainKey(ctx, g, omd1.Member())
	assert.NoError(t, err)

	ds1For2Encrypted, err := mkh1.GetShareableChainKey(ctx, g, omd2.Member())
	assert.NoError(t, err)

	ds2For2Encrypted, err := mkh2.GetShareableChainKey(ctx, g, omd2.Member())
	assert.NoError(t, err)

	err = mkh1.RegisterChainKey(ctx, g, omd1.Device(), ds1For1Encrypted)
	assert.NoError(t, err)

	err = mkh2.RegisterChainKey(ctx, g, omd2.Device(), ds2For2Encrypted)
	assert.NoError(t, err)

	err = mkh2.RegisterChainKey(ctx, g, omd1.Device(), ds1For2Encrypted)
	assert.NoError(t, err)

	ds1, err := mkh1.getOwnDeviceChainKeyForGroup(ctx, g)
	assert.NoError(t, err)

	initialCounter := ds1.Counter

	for i := 0; i < 1000; i++ {
		payloadRef, err := (&protocoltypes.EncryptedMessage{Plaintext: []byte("Test payload 1")}).Marshal()
		assert.NoError(t, err)
		envEncrypted, err := mkh1.SealEnvelope(ctx, g, payloadRef)
		assert.NoError(t, err)

		ds, err := mkh1.getDeviceChainKeyForGroupAndDevice(ctx, gPK, gc1DevicePubKey)
		if !assert.NoError(t, err) {
			t.Fatalf("failed at i = %d", i)
		}
		assert.Equal(t, ds.Counter, initialCounter+uint64(i+1))

		env, headers, err := ((*secretStore)(nil)).OpenEnvelopeHeaders(envEncrypted, g)
		if !assert.NoError(t, err) {
			t.Fatalf("failed at i = %d", i)
		}

		payloadClr, err := mkh2.OpenEnvelopePayload(ctx, env, headers, gPK, omd2.Device(), cid.Undef)
		if !assert.NoError(t, err) {
			t.Fatalf("failed at i = %d", i)
		}

		if assert.NotNil(t, headers) && assert.NotNil(t, payloadClr) {
			devRaw, err := omd1.Device().Raw()
			assert.NoError(t, err)

			assert.Equal(t, headers.DevicePK, devRaw)

			payloadClrBytes, err := payloadClr.Marshal()
			assert.NoError(t, err)
			assert.Equal(t, payloadRef, payloadClrBytes)
		} else {
			break
		}
	}
}

func mustDeviceChainKey(t testing.TB) func(ds *protocoltypes.DeviceChainKey, err error) *protocoltypes.DeviceChainKey {
	return func(ds *protocoltypes.DeviceChainKey, err error) *protocoltypes.DeviceChainKey {
		t.Helper()

		if err != nil {
			t.Fatal(err)
		}

		return ds
	}
}

func mustMessageHeaders(t testing.TB, omd OwnMemberDevice, counter uint64, payload []byte) *protocoltypes.MessageHeaders {
	t.Helper()

	pkB, err := omd.Device().Raw()
	if err != nil {
		t.Fatal(err)
	}

	sig, err := omd.DeviceSign(payload)
	if err != nil {
		t.Fatal(err)
	}

	return &protocoltypes.MessageHeaders{
		Counter:  counter,
		DevicePK: pkB,
		Sig:      sig,
	}
}

func Test_EncryptMessagePayload(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, _, err := protocoltypes.NewGroupMultiMember()
	assert.NoError(t, err)

	mkh1, err := newInMemSecretStore(nil)
	assert.NoError(t, err)
	t.Cleanup(func() {
		mkh1.Close()
	})

	mkh2, err := newInMemSecretStore(nil)
	assert.NoError(t, err)
	t.Cleanup(func() {
		mkh2.Close()
	})

	omd1, err := mkh1.GetOwnMemberDeviceForGroup(group)
	assert.NoError(t, err)

	omd2, err := mkh2.GetOwnMemberDeviceForGroup(group)
	assert.NoError(t, err)

	encryptedDS1For1, err := mkh1.GetShareableChainKey(ctx, group, omd1.Member())
	assert.NoError(t, err)

	encryptedDS1For2, err := mkh1.GetShareableChainKey(ctx, group, omd2.Member())
	assert.NoError(t, err)

	encryptedDS2For2, err := mkh2.GetShareableChainKey(ctx, group, omd2.Member())
	assert.NoError(t, err)

	gc1DevicePubKey := omd1.Device()
	gc2DevicePubKey := omd2.Device()

	err = mkh1.RegisterChainKey(ctx, group, gc1DevicePubKey, encryptedDS1For1)
	assert.NoError(t, err)

	err = mkh2.RegisterChainKey(ctx, group, gc2DevicePubKey, encryptedDS2For2)
	assert.NoError(t, err)

	ds1, err := mkh1.getOwnDeviceChainKeyForGroup(ctx, group)
	assert.NoError(t, err)

	initialCounter := ds1.Counter
	firstDeviceChainKey := append([]byte(nil), ds1.ChainKey...)

	payloadRef1 := []byte("ok, this is the first test")
	payloadRef2 := []byte("so, this is a second test")
	payloadRef3 := []byte("this will be posted many times")

	err = mkh2.RegisterChainKey(ctx, group, omd1.Device(), encryptedDS1For2)
	assert.NoError(t, err)

	gPK, err := group.GetPubKey()
	assert.NoError(t, err)

	assert.Equal(t, mustDeviceChainKey(t)(mkh1.getDeviceChainKeyForGroupAndDevice(ctx, gPK, omd1.Device())).ChainKey, firstDeviceChainKey)

	payloadEnc1, _, err := sealPayload(payloadRef1, mustDeviceChainKey(t)(mkh1.getDeviceChainKeyForGroupAndDevice(ctx, gPK, omd1.Device())), omd1.(*ownMemberDevice).device, group)
	assert.NoError(t, err)

	// secret is derived by sealEnvelope
	err = mkh1.deriveDeviceChainKey(ctx, group, omd1.Device())
	assert.NoError(t, err)

	assert.NotEqual(t, hex.EncodeToString(payloadRef1), hex.EncodeToString(payloadEnc1))

	// Messages are encrypted with DeviceChainKey.Counter
	// uint64 overflows to 0, which is the expected behaviour

	// Test with a wrong counter value
	payloadClr1, decryptInfo, err := mkh2.openPayload(ctx, cid.Undef, gPK, payloadEnc1, mustMessageHeaders(t, omd1, initialCounter+2, payloadRef1))
	assert.Error(t, err)
	assert.Nil(t, decryptInfo)
	assert.Equal(t, "", string(payloadClr1))

	// Test with a valid counter value, but no CID (so no cache)
	payloadClr1, decryptInfo, err = mkh2.openPayload(ctx, cid.Undef, gPK, payloadEnc1, mustMessageHeaders(t, omd1, initialCounter+1, payloadRef1))
	assert.NoError(t, err)
	assert.Equal(t, string(payloadRef1), string(payloadClr1))

	err = mkh2.postDecryptActions(ctx, decryptInfo, gPK, omd2.Device(), mustMessageHeaders(t, omd1, initialCounter+1, payloadRef1))
	assert.NoError(t, err)

	ds, err := mkh1.getDeviceChainKeyForGroupAndDevice(ctx, gPK, omd1.Device())
	assert.NoError(t, err)

	assert.Equal(t, ds.Counter, initialCounter+1)
	assert.NotEqual(t, ds.ChainKey, firstDeviceChainKey)

	payloadEnc2, _, err := sealPayload(payloadRef1, ds, omd1.(*ownMemberDevice).device, group)
	assert.NoError(t, err)

	err = mkh1.deriveDeviceChainKey(ctx, group, omd1.Device())
	assert.NoError(t, err)

	// Ensure that encrypted message is not the same as the first message
	assert.NotEqual(t, hex.EncodeToString(payloadRef1), hex.EncodeToString(payloadEnc2))
	assert.NotEqual(t, hex.EncodeToString(payloadEnc1), hex.EncodeToString(payloadEnc2))

	payloadClr2, decryptInfo, err := mkh2.openPayload(ctx, cid.Undef, gPK, payloadEnc2, mustMessageHeaders(t, omd1, initialCounter+2, payloadRef1))
	assert.NoError(t, err)

	err = mkh2.postDecryptActions(ctx, decryptInfo, gPK, omd2.Device(), mustMessageHeaders(t, omd1, initialCounter+2, payloadRef1))
	assert.NoError(t, err)

	assert.Equal(t, string(payloadRef1), string(payloadClr2))

	// Make sure that a message without a CID can't be decrypted twice
	payloadClr2, decryptInfo, err = mkh2.openPayload(ctx, cid.Undef, gPK, payloadEnc2, mustMessageHeaders(t, omd1, initialCounter+1, payloadRef1))
	assert.Error(t, err)
	assert.Equal(t, "", string(payloadClr2))

	ds, err = mkh1.getDeviceChainKeyForGroupAndDevice(ctx, gPK, omd1.Device())
	assert.NoError(t, err)

	// Make sure that a message a CID can be decrypted twice
	payloadEnc3, _, err := sealPayload(payloadRef2, ds, omd1.(*ownMemberDevice).device, group)
	assert.NoError(t, err)

	err = mkh1.deriveDeviceChainKey(ctx, group, omd1.Device())
	assert.NoError(t, err)

	dummyCID1, err := cid.Parse("QmbdQXQh9B2bWZgZJqfbjNPV5jGN2owbQ3vjeYsaDaCDqU")
	assert.NoError(t, err)

	dummyCID2, err := cid.Parse("Qmf8oj9wbfu73prNAA1cRQVDqA52gD5B3ApnYQQjcjffH4")
	assert.NoError(t, err)

	// Not decrypted message yet, wrong counter value
	payloadClr3, decryptInfo, err := mkh2.openPayload(ctx, dummyCID1, gPK, payloadEnc3, mustMessageHeaders(t, omd1, initialCounter+2, payloadRef2))
	assert.Error(t, err)
	assert.Equal(t, "", string(payloadClr3))

	payloadClr3, decryptInfo, err = mkh2.openPayload(ctx, dummyCID1, gPK, payloadEnc3, mustMessageHeaders(t, omd1, initialCounter+3, payloadRef2))
	assert.NoError(t, err)
	assert.Equal(t, string(payloadRef2), string(payloadClr3))

	err = mkh2.postDecryptActions(ctx, decryptInfo, gPK, omd2.Device(), mustMessageHeaders(t, omd1, initialCounter+3, payloadRef2))
	assert.NoError(t, err)

	payloadClr3, decryptInfo, err = mkh2.openPayload(ctx, dummyCID1, gPK, payloadEnc3, mustMessageHeaders(t, omd1, initialCounter+3, payloadRef2))
	assert.NoError(t, err)
	assert.Equal(t, string(payloadRef2), string(payloadClr3))

	err = mkh2.postDecryptActions(ctx, decryptInfo, gPK, omd2.Device(), mustMessageHeaders(t, omd1, initialCounter+3, payloadRef2))
	assert.NoError(t, err)

	// Wrong CID
	payloadClr3, decryptInfo, err = mkh2.openPayload(ctx, dummyCID2, gPK, payloadEnc3, mustMessageHeaders(t, omd1, initialCounter+3, payloadRef2))
	assert.Error(t, err)
	assert.Equal(t, "", string(payloadClr3))

	// Reused CID, wrong counter value
	payloadClr3, decryptInfo, err = mkh2.openPayload(ctx, dummyCID1, gPK, payloadEnc3, mustMessageHeaders(t, omd1, initialCounter+4, payloadRef2))
	assert.Error(t, err)
	assert.Equal(t, "", string(payloadClr3))

	massExpected := uint64(200)

	// Test appending 200 messages, to ensure new secrets are generated correctly
	for i := uint64(0); i < massExpected; i++ {
		ds, err = mkh1.getDeviceChainKeyForGroupAndDevice(ctx, gPK, omd1.Device())
		assert.NoError(t, err)

		payloadEnc, _, err := sealPayload(payloadRef3, ds, omd1.(*ownMemberDevice).device, group)
		assert.NoError(t, err)

		err = mkh1.deriveDeviceChainKey(ctx, group, omd1.Device())
		assert.NoError(t, err)

		ds, err = mkh1.getDeviceChainKeyForGroupAndDevice(ctx, gPK, omd1.Device())
		assert.NoError(t, err)

		counter := ds.Counter

		payloadClr, decryptInfo, err := mkh2.openPayload(ctx, cid.Undef, gPK, payloadEnc, mustMessageHeaders(t, omd1, counter, payloadRef3))
		if !assert.NoError(t, err) {
			t.Fatalf("failed at i = %d", i)
		}

		err = mkh2.postDecryptActions(ctx, decryptInfo, gPK, omd2.Device(), mustMessageHeaders(t, omd1, counter, payloadRef3))
		assert.NoError(t, err)

		assert.Equal(t, string(payloadRef3), string(payloadClr))
	}

	ds, err = mkh1.getDeviceChainKeyForGroupAndDevice(ctx, gPK, omd1.Device())
	assert.NoError(t, err)

	assert.Equal(t, initialCounter+massExpected+3, ds.Counter)
}

func TestGetGroupForContact(t *testing.T) {
	privateKey, _, err := crypto.GenerateEd25519Key(crand.Reader)
	require.NoError(t, err)

	group, err := getGroupForContact(privateKey)
	require.NoError(t, err)

	require.Equal(t, group.GroupType, protocoltypes.GroupTypeContact)
	require.Equal(t, len(group.PublicKey), 32)
	require.Equal(t, len(group.Secret), 32)
}

func TestGetKeysForGroupOfContact(t *testing.T) {
	privateKey, _, err := crypto.GenerateEd25519Key(crand.Reader)
	require.NoError(t, err)

	groupPrivateKey, groupSecretPrivateKey, err := getKeysForGroupOfContact(privateKey)
	require.NoError(t, err)

	require.NotNil(t, groupPrivateKey)
	require.NotNil(t, groupSecretPrivateKey)
	require.False(t, groupPrivateKey.Equals(groupSecretPrivateKey))
}
