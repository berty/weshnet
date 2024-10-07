package weshnet

import (
	"context"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"berty.tech/weshnet/v2/pkg/protocoltypes"
	"berty.tech/weshnet/v2/pkg/secretstore"
	"berty.tech/weshnet/v2/pkg/testutil"
)

func Test_sealPushMessage_OutOfStoreReceive(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp, cancel := NewTestingProtocol(ctx, t, &TestingOpts{}, nil)
	defer cancel()

	g, _, err := NewGroupMultiMember()
	require.NoError(t, err)

	s := tp.Service

	gPK, err := g.GetPubKey()
	require.NoError(t, err)

	_, err = s.MultiMemberGroupJoin(ctx, &protocoltypes.MultiMemberGroupJoin_Request{Group: g})
	require.NoError(t, err)

	gPKRaw, err := gPK.Raw()
	require.NoError(t, err)

	_, err = s.ActivateGroup(ctx, &protocoltypes.ActivateGroup_Request{GroupPk: gPKRaw})
	require.NoError(t, err)

	gc, err := s.(ServiceMethods).GetContextGroupForID(g.PublicKey)
	require.NoError(t, err)

	otherSecretStore, cancel := createVirtualOtherPeerSecrets(t, ctx, gc)
	defer cancel()

	testPayload := []byte("test payload")

	envBytes, err := otherSecretStore.SealEnvelope(ctx, g, testPayload)
	require.NoError(t, err)

	env, headers, err := otherSecretStore.OpenEnvelopeHeaders(envBytes, g)
	require.NoError(t, err)

	oosMsgEnv, err := otherSecretStore.SealOutOfStoreMessageEnvelope(cid.Undef, env, headers, g)
	require.NoError(t, err)
	oosMsgEnvBytes, err := proto.Marshal(oosMsgEnv)
	require.NoError(t, err)

	outOfStoreMessage, group, clearPayload, alreadyDecrypted, err := gc.SecretStore().OpenOutOfStoreMessage(ctx, oosMsgEnvBytes)
	require.NoError(t, err)

	require.Equal(t, g.PublicKey, group.PublicKey)
	require.Equal(t, g.Secret, group.Secret)
	require.Equal(t, g.SecretSig, group.SecretSig)
	require.Equal(t, g.GroupType, group.GroupType)
	require.Equal(t, g.SignPub, group.SignPub)
	require.Equal(t, g.LinkKey, group.LinkKey)
	require.Equal(t, g.LinkKeySig, group.LinkKeySig)
	require.Equal(t, []byte("test payload"), clearPayload)
	require.False(t, alreadyDecrypted)

	require.Equal(t, headers.Counter, outOfStoreMessage.Counter)
	require.Equal(t, headers.DevicePk, outOfStoreMessage.DevicePk)
	require.Equal(t, headers.Sig, outOfStoreMessage.Sig)
	require.Equal(t, env.Message, outOfStoreMessage.EncryptedPayload)
}

func Test_OutOfStoreMessageFlow(t *testing.T) {
	message := []byte("test message")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, cleanup := testutil.Logger(t)
	defer cleanup()

	tp, cancel := NewTestingProtocol(ctx, t, &TestingOpts{Logger: logger}, nil)
	defer cancel()

	g, _, err := NewGroupMultiMember()
	require.NoError(t, err)

	s := tp.Service

	gPK, err := g.GetPubKey()
	require.NoError(t, err)

	_, err = s.MultiMemberGroupJoin(ctx, &protocoltypes.MultiMemberGroupJoin_Request{Group: g})
	require.NoError(t, err)

	gPKRaw, err := gPK.Raw()
	require.NoError(t, err)

	_, err = s.ActivateGroup(ctx, &protocoltypes.ActivateGroup_Request{GroupPk: gPKRaw})
	require.NoError(t, err)

	// send a message
	sendReply, err := s.AppMessageSend(ctx, &protocoltypes.AppMessageSend_Request{
		GroupPk: gPKRaw,
		Payload: message,
	})
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// craft an out of store message
	craftReply, err := s.OutOfStoreSeal(ctx, &protocoltypes.OutOfStoreSeal_Request{
		Cid:            sendReply.Cid,
		GroupPublicKey: gPKRaw,
	})
	require.NoError(t, err)

	// verify the out of store message
	openReply, err := s.OutOfStoreReceive(ctx, &protocoltypes.OutOfStoreReceive_Request{
		Payload: craftReply.Encrypted,
	})
	require.NoError(t, err)

	encryptedMessage := &protocoltypes.EncryptedMessage{}
	err = proto.Unmarshal(openReply.Cleartext, encryptedMessage)
	require.NoError(t, err)

	require.Equal(t, message, encryptedMessage.Plaintext)
}

func createVirtualOtherPeerSecrets(t testing.TB, ctx context.Context, gc *GroupContext) (secretstore.SecretStore, func()) {
	secretStore, err := secretstore.NewInMemSecretStore(nil)
	require.NoError(t, err)

	cleanup := func() {
		_ = secretStore.Close()
	}

	// Manually adding another member to the group
	otherMD, err := secretStore.GetOwnMemberDeviceForGroup(gc.Group())
	require.NoError(t, err)
	_, err = MetadataStoreAddDeviceToGroup(ctx, gc.MetadataStore(), gc.Group(), otherMD)
	require.NoError(t, err)

	memberDevice, err := gc.SecretStore().GetOwnMemberDeviceForGroup(gc.Group())
	require.NoError(t, err)

	ds, err := secretStore.GetShareableChainKey(ctx, gc.Group(), memberDevice.Member())
	require.NoError(t, err)

	_, err = MetadataStoreSendSecret(ctx, gc.MetadataStore(), gc.Group(), otherMD, memberDevice.Member(), ds)
	require.NoError(t, err)

	time.Sleep(time.Millisecond * 200)

	return secretStore, cleanup
}
