package weshnet

import (
	"context"
	"io"
	"testing"
	"time"

	libp2p_mocknet "github.com/berty/go-libp2p-mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/testutil"
)

func TestContactRequestFlow(t *testing.T) {
	testutil.FilterSpeed(t, testutil.Slow)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	logger, cleanup := testutil.Logger(t)
	defer cleanup()

	opts := TestingOpts{
		Mocknet: libp2p_mocknet.New(),
		Logger:  logger,
	}

	metadataSender1 := []byte("sender_1")

	pts, cleanup := NewTestingProtocolWithMockedPeers(ctx, t, &opts, nil, 2)
	defer cleanup()

	_, err := pts[0].Client.ContactRequestEnable(ctx, &protocoltypes.ContactRequestEnable_Request{})
	require.NoError(t, err)

	_, err = pts[1].Client.ContactRequestEnable(ctx, &protocoltypes.ContactRequestEnable_Request{})
	require.NoError(t, err)

	config0, err := pts[0].Client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
	require.NoError(t, err)
	require.NotNil(t, config0)

	config1, err := pts[1].Client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
	require.NoError(t, err)
	require.NotNil(t, config1)

	ref0, err := pts[0].Client.ContactRequestResetReference(ctx, &protocoltypes.ContactRequestResetReference_Request{})
	require.NoError(t, err)
	require.NotNil(t, ref0)

	ref1, err := pts[1].Client.ContactRequestResetReference(ctx, &protocoltypes.ContactRequestResetReference_Request{})
	require.NoError(t, err)
	require.NotNil(t, ref1)

	subCtx, subCancel := context.WithCancel(ctx)
	defer subCancel()

	subMeta0, err := pts[0].Client.GroupMetadataList(subCtx, &protocoltypes.GroupMetadataList_Request{
		GroupPk: config0.AccountGroupPk,
	})
	require.NoError(t, err)
	found := false

	_, err = pts[1].Client.ContactRequestSend(ctx, &protocoltypes.ContactRequestSend_Request{
		Contact: &protocoltypes.ShareableContact{
			Pk:                   config0.AccountPk,
			PublicRendezvousSeed: ref0.PublicRendezvousSeed,
		},
		OwnMetadata: metadataSender1,
	})
	require.NoError(t, err)

	for {
		evt, err := subMeta0.Recv()
		if err == io.EOF || subMeta0.Context().Err() != nil {
			break
		}

		require.NoError(t, err)

		if evt == nil || evt.Metadata.EventType != protocoltypes.EventType_EventTypeAccountContactRequestIncomingReceived {
			continue
		}

		req := &protocoltypes.AccountContactRequestIncomingReceived{}
		err = proto.Unmarshal(evt.Event, req)

		require.NoError(t, err)
		require.Equal(t, config1.AccountPk, req.ContactPk)
		require.Equal(t, metadataSender1, req.ContactMetadata)
		found = true
		subCancel()
	}

	require.True(t, found)

	_, err = pts[1].Client.ContactRequestAccept(ctx, &protocoltypes.ContactRequestAccept_Request{
		ContactPk: config0.AccountPk,
	})

	require.Error(t, err)

	_, err = pts[1].Client.ContactRequestAccept(ctx, &protocoltypes.ContactRequestAccept_Request{
		ContactPk: config1.AccountPk,
	})

	require.Error(t, err)

	_, err = pts[0].Client.ContactRequestAccept(ctx, &protocoltypes.ContactRequestAccept_Request{
		ContactPk: config0.AccountPk,
	})

	require.Error(t, err)

	_, err = pts[0].Client.ContactRequestAccept(ctx, &protocoltypes.ContactRequestAccept_Request{
		ContactPk: config1.AccountPk,
	})

	require.NoError(t, err)

	grpInfo, err := pts[0].Client.GroupInfo(ctx, &protocoltypes.GroupInfo_Request{
		ContactPk: config1.AccountPk,
	})
	require.NoError(t, err)

	_, err = pts[0].Client.ActivateGroup(ctx, &protocoltypes.ActivateGroup_Request{
		GroupPk: grpInfo.Group.PublicKey,
	})

	require.NoError(t, err)

	_, err = pts[1].Client.ActivateGroup(ctx, &protocoltypes.ActivateGroup_Request{
		GroupPk: grpInfo.Group.PublicKey,
	})

	require.NoError(t, err)
}

func TestContactRequestFlowWithoutIncoming(t *testing.T) {
	t.Skip("KUBO: this test timeout, disable it for now")

	testutil.FilterSpeed(t, testutil.Slow)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger, cleanup := testutil.Logger(t)
	defer cleanup()

	mn := libp2p_mocknet.New()
	defer mn.Close()

	opts := TestingOpts{
		Mocknet: mn,
		Logger:  logger,
	}

	metadataSender1 := []byte("sender_1")

	pts, cleanup := NewTestingProtocolWithMockedPeers(ctx, t, &opts, nil, 2)
	defer cleanup()

	_, err := pts[0].Client.ContactRequestEnable(ctx, &protocoltypes.ContactRequestEnable_Request{})
	require.NoError(t, err)

	config0, err := pts[0].Client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
	require.NoError(t, err)
	require.NotNil(t, config0)

	config1, err := pts[1].Client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
	require.NoError(t, err)
	require.NotNil(t, config1)

	ref0, err := pts[0].Client.ContactRequestResetReference(ctx, &protocoltypes.ContactRequestResetReference_Request{})
	require.NoError(t, err)
	require.NotNil(t, ref0)

	subCtx, subCancel := context.WithCancel(ctx)
	defer subCancel()

	subMeta0, err := pts[0].Client.GroupMetadataList(subCtx, &protocoltypes.GroupMetadataList_Request{
		GroupPk: config0.AccountGroupPk,
	})
	require.NoError(t, err)
	found := false

	_, err = pts[1].Client.ContactRequestSend(ctx, &protocoltypes.ContactRequestSend_Request{
		Contact: &protocoltypes.ShareableContact{
			Pk:                   config0.AccountPk,
			PublicRendezvousSeed: ref0.PublicRendezvousSeed,
		},
		OwnMetadata: metadataSender1,
	})
	require.NoError(t, err)

	for {
		evt, err := subMeta0.Recv()
		if err != nil {
			assert.NoError(t, err)
			break
		}

		require.NoError(t, err)

		if evt == nil || evt.Metadata.EventType != protocoltypes.EventType_EventTypeAccountContactRequestIncomingReceived {
			continue
		}

		req := &protocoltypes.AccountContactRequestIncomingReceived{}
		err = proto.Unmarshal(evt.Event, req)

		require.NoError(t, err)
		require.Equal(t, config1.AccountPk, req.ContactPk)
		require.Equal(t, metadataSender1, req.ContactMetadata)
		found = true
		subCancel()
	}

	require.True(t, found)

	_, err = pts[0].Client.ContactRequestAccept(ctx, &protocoltypes.ContactRequestAccept_Request{
		ContactPk: config1.AccountPk,
	})

	require.NoError(t, err)

	_, err = pts[0].Client.GroupInfo(ctx, &protocoltypes.GroupInfo_Request{
		ContactPk: config1.AccountPk,
	})
	require.NoError(t, err)
}
