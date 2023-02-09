package weshnet_test

import (
	"context"
	crand "crypto/rand"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/nacl/box"

	"berty.tech/berty/v2/go/pkg/bertypush"
	weshnet "berty.tech/weshnet"
	"berty.tech/weshnet/pkg/cryptoutil"
	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/pushtypes"
)

func Test_sealPushMessage_decryptOutOfStoreMessageEnv(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, devicePushSK, err := box.GenerateKey(crand.Reader)
	require.NoError(t, err)

	tp, cancel := weshnet.NewTestingProtocol(ctx, t, &weshnet.TestingOpts{PushSK: devicePushSK}, nil)
	defer cancel()

	g, _, err := weshnet.NewGroupMultiMember()
	require.NoError(t, err)

	s := tp.Service

	gPK, err := g.GetPubKey()
	require.NoError(t, err)

	_, err = s.MultiMemberGroupJoin(ctx, &protocoltypes.MultiMemberGroupJoin_Request{Group: g})
	require.NoError(t, err)

	gPKRaw, err := gPK.Raw()
	require.NoError(t, err)

	_, err = s.ActivateGroup(ctx, &protocoltypes.ActivateGroup_Request{GroupPK: gPKRaw})
	require.NoError(t, err)

	gc, err := s.(weshnet.ServiceMethods).GetContextGroupForID(g.PublicKey)
	require.NoError(t, err)

	otherMD, otherDS := weshnet.CreateVirtualOtherPeerSecretsShareSecret(t, ctx, []*weshnet.MetadataStore{gc.MetadataStore()})

	testPayload := []byte("test payload")

	envBytes, err := cryptoutil.SealEnvelope(testPayload, otherDS, otherMD.PrivateDevice(), g)
	require.NoError(t, err)

	env, headers, err := cryptoutil.OpenEnvelopeHeaders(envBytes, g)
	require.NoError(t, err)

	oosMsgEnv, err := weshnet.SealOutOfStoreMessageEnvelope(cid.Undef, env, headers, g)
	require.NoError(t, err)

	openedOOSMessage, err := bertypush.DecryptOutOfStoreMessageEnv(ctx, tp.GroupDatastore, oosMsgEnv, gPK)
	require.NoError(t, err)

	require.Equal(t, headers.Counter, openedOOSMessage.Counter)
	require.Equal(t, headers.DevicePK, openedOOSMessage.DevicePK)
	require.Equal(t, headers.Sig, openedOOSMessage.Sig)
	require.Equal(t, env.Message, openedOOSMessage.EncryptedPayload)
}

func TestService_PushShareToken(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	devicePushPK, devicePushSK, err := box.GenerateKey(crand.Reader)
	require.NoError(t, err)

	serverPushPK, _, err := box.GenerateKey(crand.Reader)
	require.NoError(t, err)

	tokenTestData := []byte("token_test_data_1")
	const nameTestPackage = "test.app"
	const serverAddr1 = "server1.test"

	tp, cancel := weshnet.NewTestingProtocol(ctx, t, &weshnet.TestingOpts{PushSK: devicePushSK}, nil)
	defer cancel()

	s := tp.Service

	_, err = s.PushSetServer(ctx, &protocoltypes.PushSetServer_Request{Server: &protocoltypes.PushServer{
		ServerKey:   serverPushPK[:],
		ServiceAddr: serverAddr1,
	}})
	require.NoError(t, err)

	_, err = s.PushSetDeviceToken(ctx, &protocoltypes.PushSetDeviceToken_Request{
		Receiver: &protocoltypes.PushServiceReceiver{
			TokenType:          pushtypes.PushServiceTokenType_PushTokenApplePushNotificationService,
			BundleID:           nameTestPackage,
			Token:              tokenTestData,
			RecipientPublicKey: devicePushPK[:],
		},
	})
	require.NoError(t, err)

	g, gSK, err := weshnet.NewGroupMultiMember()
	require.NoError(t, err)

	_, err = s.MultiMemberGroupJoin(ctx, &protocoltypes.MultiMemberGroupJoin_Request{Group: g})
	require.NoError(t, err)

	gPK, _ := gSK.GetPublic().Raw()
	require.NoError(t, err)

	_, err = s.ActivateGroup(ctx, &protocoltypes.ActivateGroup_Request{GroupPK: gPK})
	require.NoError(t, err)

	gc, err := s.(weshnet.ServiceMethods).GetContextGroupForID(g.PublicKey)
	require.NoError(t, err)

	pushToken, err := gc.MetadataStore().GetPushTokenForDevice(gc.DevicePubKey())
	require.Error(t, err)
	require.Nil(t, pushToken)

	_, err = s.PushShareToken(ctx, &protocoltypes.PushShareToken_Request{
		GroupPK: g.PublicKey,
		Server: &protocoltypes.PushServer{
			ServerKey:   serverPushPK[:],
			ServiceAddr: serverAddr1,
		},
		Receiver: &protocoltypes.PushServiceReceiver{
			TokenType:          pushtypes.PushServiceTokenType_PushTokenApplePushNotificationService,
			BundleID:           nameTestPackage,
			Token:              tokenTestData,
			RecipientPublicKey: devicePushPK[:],
		},
	})
	require.NoError(t, err)

	pushToken, err = gc.MetadataStore().GetPushTokenForDevice(gc.DevicePubKey())
	require.NoError(t, err)
	require.NotNil(t, pushToken)
	require.Equal(t, serverAddr1, pushToken.Server.ServiceAddr)
	require.Equal(t, serverPushPK[:], pushToken.Server.ServerKey)
}

func TestService_PushSetDeviceToken(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	devicePushPK, devicePushSK, err := box.GenerateKey(crand.Reader)
	require.NoError(t, err)

	tokenTestData1 := []byte("token_test_data_1")
	tokenTestData2 := []byte("token_test_data_2")
	const nameTestPackage = "test.app"

	tp, cancel := weshnet.NewTestingProtocol(ctx, t, &weshnet.TestingOpts{PushSK: devicePushSK}, nil)
	defer cancel()

	s := tp.Service

	currentPush, _ := s.(weshnet.ServiceMethods).GetCurrentDevicePushConfig()
	require.Nil(t, currentPush)

	_, err = s.PushSetDeviceToken(ctx, &protocoltypes.PushSetDeviceToken_Request{
		Receiver: &protocoltypes.PushServiceReceiver{
			TokenType: pushtypes.PushServiceTokenType_PushTokenMQTT,
			BundleID:  nameTestPackage,
			Token:     tokenTestData1,
		},
	})
	require.NoError(t, err)

	currentPush, _ = s.(weshnet.ServiceMethods).GetCurrentDevicePushConfig()
	require.NotNil(t, currentPush)
	require.Equal(t, tokenTestData1, currentPush.Token)
	require.Equal(t, nameTestPackage, currentPush.BundleID)
	require.Equal(t, pushtypes.PushServiceTokenType_PushTokenMQTT, currentPush.TokenType)
	require.Equal(t, devicePushPK[:], currentPush.RecipientPublicKey)

	_, err = s.PushSetDeviceToken(ctx, &protocoltypes.PushSetDeviceToken_Request{
		Receiver: &protocoltypes.PushServiceReceiver{
			TokenType: pushtypes.PushServiceTokenType_PushTokenApplePushNotificationService,
			BundleID:  nameTestPackage,
			Token:     tokenTestData2,
		},
	})
	require.NoError(t, err)

	currentPush, _ = s.(weshnet.ServiceMethods).GetCurrentDevicePushConfig()
	require.NotNil(t, currentPush)
	require.Equal(t, tokenTestData2, currentPush.Token)
	require.Equal(t, nameTestPackage, currentPush.BundleID)
	require.Equal(t, pushtypes.PushServiceTokenType_PushTokenApplePushNotificationService, currentPush.TokenType)
	require.Equal(t, devicePushPK[:], currentPush.RecipientPublicKey)
}

func TestService_PushSetServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const serverAddr1 = "server1.test"
	const serverAddr2 = "server2.test"

	_, devicePushSK, err := box.GenerateKey(crand.Reader)
	require.NoError(t, err)

	serverPushPK1, _, err := box.GenerateKey(crand.Reader)
	require.NoError(t, err)

	serverPushPK2, _, err := box.GenerateKey(crand.Reader)
	require.NoError(t, err)

	tp, cancel := weshnet.NewTestingProtocol(ctx, t, &weshnet.TestingOpts{PushSK: devicePushSK}, nil)
	defer cancel()

	s := tp.Service

	_, currentPush := s.(weshnet.ServiceMethods).GetCurrentDevicePushConfig()
	require.Nil(t, currentPush)

	_, err = s.PushSetServer(ctx, &protocoltypes.PushSetServer_Request{Server: &protocoltypes.PushServer{
		ServerKey:   serverPushPK1[:],
		ServiceAddr: serverAddr1,
	}})
	require.NoError(t, err)

	_, currentPush = s.(weshnet.ServiceMethods).GetCurrentDevicePushConfig()
	require.NotNil(t, currentPush)
	require.Equal(t, serverPushPK1[:], currentPush.ServerKey)
	require.Equal(t, serverAddr1, currentPush.ServiceAddr)

	_, err = s.PushSetServer(ctx, &protocoltypes.PushSetServer_Request{Server: &protocoltypes.PushServer{
		ServerKey:   serverPushPK2[:],
		ServiceAddr: serverAddr2,
	}})
	require.NoError(t, err)

	_, currentPush = s.(weshnet.ServiceMethods).GetCurrentDevicePushConfig()
	require.NotNil(t, currentPush)
	require.Equal(t, serverPushPK2[:], currentPush.ServerKey)
	require.Equal(t, serverAddr2, currentPush.ServiceAddr)

	// FIXME: Should we add a way to clear the push server used with the current device?
}
