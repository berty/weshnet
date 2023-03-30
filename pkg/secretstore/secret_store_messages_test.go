package secretstore_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	cid "github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/stretchr/testify/assert"

	"berty.tech/weshnet"
	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/secretstore"
	"berty.tech/weshnet/pkg/testutil"
)

func addDummyMemberInMetadataStore(ctx context.Context, t testing.TB, ms *weshnet.MetadataStore, g *protocoltypes.Group, memberPK crypto.PubKey, join bool) crypto.PubKey {
	t.Helper()

	secretStore, err := secretstore.NewInMemSecretStore(nil)
	assert.NoError(t, err)
	assert.NotNil(t, secretStore)

	md, err := secretStore.GetOwnMemberDeviceForGroup(g)
	assert.NoError(t, err)

	if join {
		_, err = weshnet.MetadataStoreAddDeviceToGroup(ctx, ms, g, md)
		assert.NoError(t, err)
	}

	deviceChainKeyToSend, err := secretStore.GetShareableChainKey(ctx, g, memberPK)
	assert.NoError(t, err)

	_, err = weshnet.MetadataStoreSendSecret(ctx, ms, g, md, memberPK, deviceChainKeyToSend)
	assert.NoError(t, err)

	return md.Device()
}

func Test_EncryptMessageEnvelope(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, _, err := weshnet.NewGroupMultiMember()
	assert.NoError(t, err)

	secretStore1, err := secretstore.NewInMemSecretStore(nil)
	assert.NoError(t, err)
	assert.NotNil(t, secretStore1)

	t.Cleanup(func() {
		_ = secretStore1.Close()
	})

	omd1, err := secretStore1.GetOwnMemberDeviceForGroup(g)
	assert.NoError(t, err)

	gc1 := weshnet.NewContextGroup(g, nil, nil, secretStore1, omd1, nil)

	deviceChainKey1For1, err := secretStore1.GetShareableChainKey(ctx, g, gc1.MemberPubKey())
	assert.NoError(t, err)

	err = secretStore1.RegisterChainKey(ctx, g, gc1.DevicePubKey(), deviceChainKey1For1)
	assert.NoError(t, err)

	secretStore2, err := secretstore.NewInMemSecretStore(nil)
	assert.NoError(t, err)
	assert.NotNil(t, secretStore2)

	t.Cleanup(func() {
		_ = secretStore2.Close()
	})

	omd2, err := secretStore2.GetOwnMemberDeviceForGroup(g)
	assert.NoError(t, err)

	payloadRef1, err := (&protocoltypes.EncryptedMessage{Plaintext: []byte("Test payload 1")}).Marshal()
	assert.NoError(t, err)

	deviceChainKey1For2, err := secretStore1.GetShareableChainKey(ctx, g, omd2.Member())
	assert.NoError(t, err)

	deviceChainKey2For2, err := secretStore2.GetShareableChainKey(ctx, g, omd2.Member())
	assert.NoError(t, err)

	err = secretStore2.RegisterChainKey(ctx, g, omd2.Device(), deviceChainKey2For2)
	assert.NoError(t, err)

	err = secretStore2.RegisterChainKey(ctx, g, omd1.Device(), deviceChainKey1For2)
	assert.NoError(t, err)

	env1, err := secretStore1.SealEnvelope(ctx, g, payloadRef1)
	assert.NoError(t, err)

	headers, payloadClr1, err := openEnvelope(ctx, t, secretStore2, g, omd2.Device(), env1, cid.Undef)
	assert.NoError(t, err)

	devRaw, err := omd1.Device().Raw()
	assert.Equal(t, headers.DevicePK, devRaw)

	payloadClrlBytes, err := payloadClr1.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, payloadRef1, payloadClrlBytes)
}

func testMessageKeyHolderCatchUp(t *testing.T, expectedNewDevices int, isSlow bool) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if isSlow {
		testutil.FilterSpeed(t, testutil.Slow)
	}

	dir := path.Join(os.TempDir(), fmt.Sprintf("%d", os.Getpid()), "MessageKeyHolderCatchUp")
	defer os.RemoveAll(dir)

	peers, _, cleanup := weshnet.CreatePeersWithGroupTest(ctx, t, dir, 1, 1)
	defer cleanup()

	peer := peers[0]

	secretStore1 := peer.SecretStore
	ms1 := peer.GC.MetadataStore()

	groupPublicKey, err := peer.GC.Group().GetPubKey()
	assert.NoError(t, err)

	devicesPK := make([]crypto.PubKey, expectedNewDevices)

	for i := 0; i < expectedNewDevices; i++ {
		devicesPK[i] = addDummyMemberInMetadataStore(ctx, t, ms1, peer.GC.Group(), peer.GC.MemberPubKey(), true)
	}

	for range peer.GC.FillMessageKeysHolderUsingPreviousData() {
	}

	for i, devicePublicKey := range devicesPK {
		if !assert.True(t, secretStore1.IsChainKeyKnownForDevice(ctx, groupPublicKey, devicePublicKey)) {
			t.Fatalf("failed at iteration %d", i)
		}
	}
}

func TestMessageKeyHolderCatchUp(t *testing.T) {
	for _, testCase := range []struct {
		expectedNewDevices int
		slow               bool
	}{
		{
			expectedNewDevices: 2,
			slow:               false,
		},
		{
			expectedNewDevices: 10,
			slow:               true,
		},
	} {
		testMessageKeyHolderCatchUp(t, testCase.expectedNewDevices, testCase.slow)
	}
}

func testMessageKeyHolderSubscription(t *testing.T, expectedNewDevices int, isSlow bool) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if isSlow {
		testutil.FilterSpeed(t, testutil.Slow)
	}

	dir := path.Join(os.TempDir(), fmt.Sprintf("%d", os.Getpid()), "MessageKeyHolderSubscription")
	defer os.RemoveAll(dir)

	peers, groupPrivateKey, cleanup := weshnet.CreatePeersWithGroupTest(ctx, t, dir, 1, 1)
	defer cleanup()

	peer := peers[0]

	secretStore1 := peer.SecretStore
	ms1 := peer.GC.MetadataStore()

	devicesPK := make([]crypto.PubKey, expectedNewDevices)

	ch := peer.GC.FillMessageKeysHolderUsingNewData()

	for i := 0; i < expectedNewDevices; i++ {
		devicesPK[i] = addDummyMemberInMetadataStore(ctx, t, ms1, peer.GC.Group(), peer.GC.MemberPubKey(), true)
	}

	i := 0
	for range ch {
		i++
		if i == expectedNewDevices {
			break
		}
	}

	for i, devicePublicKey := range devicesPK {
		if !assert.True(t, secretStore1.IsChainKeyKnownForDevice(ctx, groupPrivateKey.GetPublic(), devicePublicKey)) {
			t.Fatalf("failed at iteration %d", i)
		}
	}
}

func TestMessageKeyHolderSubscription(t *testing.T) {
	for _, testCase := range []struct {
		expectedNewDevices int
		slow               bool
	}{
		{
			expectedNewDevices: 2,
			slow:               false,
		},
		{
			expectedNewDevices: 10,
			slow:               true,
		},
	} {
		testMessageKeyHolderSubscription(t, testCase.expectedNewDevices, testCase.slow)
	}
}

// openEnvelope opens a MessageEnvelope and returns the decrypted message.
// It performs all the necessary steps to decrypt the message.
func openEnvelope(ctx context.Context, t testing.TB, secretStore secretstore.SecretStore, g *protocoltypes.Group, ownPK crypto.PubKey, data []byte, id cid.Cid) (*protocoltypes.MessageHeaders, *protocoltypes.EncryptedMessage, error) {
	t.Helper()

	assert.NotNil(t, secretStore)
	assert.NotNil(t, g)

	env, headers, err := secretStore.OpenEnvelopeHeaders(data, g)
	assert.NoError(t, err)

	gPK, err := g.GetPubKey()
	assert.NoError(t, err)

	msg, err := secretStore.OpenEnvelopePayload(ctx, env, headers, gPK, ownPK, id)
	assert.NoError(t, err)

	return headers, msg, nil
}
