package weshnet_test

import (
	"context"
	"testing"
	"time"

	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	libp2p_mocknet "github.com/berty/go-libp2p-mock"
	"github.com/stretchr/testify/require"

	"berty.tech/weshnet"
	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/testutil"
	"berty.tech/weshnet/pkg/tinder"
)

func TestReactivateAccountGroup(t *testing.T) {
	testutil.FilterStability(t, testutil.Flappy)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	logger, cleanup := testutil.Logger(t)
	defer cleanup()

	mn := libp2p_mocknet.New()
	defer mn.Close()

	msrv := tinder.NewMockDriverServer()

	// Setup 3 nodes
	dsA := dsync.MutexWrap(ds.NewMapDatastore())
	nodeA, closeNodeA := weshnet.NewTestingProtocol(ctx, t, &weshnet.TestingOpts{
		Logger:          logger.Named("nodeA"),
		Mocknet:         mn,
		DiscoveryServer: msrv,
	}, dsA)
	defer closeNodeA()

	dsB := dsync.MutexWrap(ds.NewMapDatastore())
	nodeB, closeNodeB := weshnet.NewTestingProtocol(ctx, t, &weshnet.TestingOpts{
		Logger:          logger.Named("nodeB"),
		Mocknet:         mn,
		DiscoveryServer: msrv,
	}, dsB)
	defer closeNodeB()

	dsC := dsync.MutexWrap(ds.NewMapDatastore())
	nodeC, closeNodeC := weshnet.NewTestingProtocol(ctx, t, &weshnet.TestingOpts{
		Logger:          logger.Named("nodeC"),
		Mocknet:         mn,
		DiscoveryServer: msrv,
	}, dsC)
	defer closeNodeC()

	// make connections
	err := mn.LinkAll()
	require.NoError(t, err)

	err = mn.ConnectAllButSelf()
	require.NoError(t, err)

	// test communication between nodeA and nodeB
	nodes := []*weshnet.TestingProtocol{nodeA, nodeB}
	addAsContact(ctx, t, nodes, nodes)
	sendMessageToContact(ctx, t, []string{"pre-deactivate nodeA-nodeB"}, nodes)

	// reactivate nodeA account group
	nodeACfg, err := nodeA.Client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
	require.NoError(t, err)
	require.NotNil(t, nodeACfg)

	_, err = nodeA.Client.DeactivateGroup(ctx, &protocoltypes.DeactivateGroup_Request{
		GroupPK: nodeACfg.AccountGroupPK,
	})
	require.NoError(t, err)

	_, err = nodeA.Client.ActivateGroup(ctx, &protocoltypes.ActivateGroup_Request{
		GroupPK: nodeACfg.AccountGroupPK,
	})
	require.NoError(t, err)

	// test communication between nodeA and nodeC
	nodes = []*weshnet.TestingProtocol{nodeA, nodeC}

	addAsContact(ctx, t, nodes, nodes)
	sendMessageToContact(ctx, t, []string{"post reactivate nodeA-nodeC"}, nodes)
}

func TestRaceReactivateAccountGroup(t *testing.T) {
	testutil.FilterStability(t, testutil.Flappy)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	logger, cleanup := testutil.Logger(t)
	defer cleanup()

	mn := libp2p_mocknet.New()
	defer mn.Close()

	msrv := tinder.NewMockDriverServer()

	// Setup 2 nodes
	dsA := dsync.MutexWrap(ds.NewMapDatastore())
	nodeA, closeNodeA := weshnet.NewTestingProtocol(ctx, t, &weshnet.TestingOpts{
		Logger:          logger.Named("nodeA"),
		Mocknet:         mn,
		DiscoveryServer: msrv,
	}, dsA)
	defer closeNodeA()

	dsB := dsync.MutexWrap(ds.NewMapDatastore())
	nodeB, closeNodeB := weshnet.NewTestingProtocol(ctx, t, &weshnet.TestingOpts{
		Logger:          logger.Named("nodeB"),
		Mocknet:         mn,
		DiscoveryServer: msrv,
	}, dsB)
	defer closeNodeB()

	// make connections
	err := mn.LinkAll()
	require.NoError(t, err)

	err = mn.ConnectAllButSelf()
	require.NoError(t, err)

	// reactivate nodeA account group
	nodeACfg, err := nodeA.Client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
	require.NoError(t, err)
	require.NotNil(t, nodeACfg)

	deactivateFunc := func() {
		t.Log("DeactivateGroup")
		_, err := nodeA.Client.DeactivateGroup(ctx, &protocoltypes.DeactivateGroup_Request{
			GroupPK: nodeACfg.AccountGroupPK,
		})
		require.NoError(t, err)
	}

	activateFunc := func() {
		t.Log("ActivateGroup")
		_, err := nodeA.Client.ActivateGroup(ctx, &protocoltypes.ActivateGroup_Request{
			GroupPK: nodeACfg.AccountGroupPK,
		})
		require.NoError(t, err)
	}

	go deactivateFunc()
	time.Sleep(1 * time.Millisecond)
	go activateFunc()

	// test communication between nodeA and nodeB
	time.Sleep(3 * time.Second)
	nodes := []*weshnet.TestingProtocol{nodeA, nodeB}
	t.Log("addAsContact")
	addAsContact(ctx, t, nodes, nodes)
	t.Log("sendMessageToContact")
	sendMessageToContact(ctx, t, []string{"nodeA-nodeB"}, nodes)
}

func TestReactivateContactGroup(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	logger, cleanup := testutil.Logger(t)
	defer cleanup()

	opts := weshnet.TestingOpts{
		Mocknet:     libp2p_mocknet.New(),
		Logger:      logger,
		ConnectFunc: weshnet.ConnectAll,
	}

	nodes, cleanup := weshnet.NewTestingProtocolWithMockedPeers(ctx, t, &opts, nil, 2)
	defer cleanup()

	addAsContact(ctx, t, nodes, nodes)

	// send messages before deactivating
	sendMessageToContact(ctx, t, []string{"pre-deactivate"}, nodes)

	// get contact group
	contactGroup := getContactGroup(ctx, t, nodes[0], nodes[1])

	// deactivate contact group
	_, err := nodes[0].Client.DeactivateGroup(ctx, &protocoltypes.DeactivateGroup_Request{
		GroupPK: contactGroup.Group.PublicKey,
	})
	require.NoError(t, err)

	// reactivate group
	_, err = nodes[0].Client.ActivateGroup(ctx, &protocoltypes.ActivateGroup_Request{
		GroupPK: contactGroup.Group.PublicKey,
	})
	require.NoError(t, err)

	// send message after reactivating
	sendMessageToContact(ctx, t, []string{"post-reactivate"}, nodes)
}

func TestRaceReactivateContactGroup(t *testing.T) {
	testutil.FilterStability(t, testutil.Flappy)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	logger, cleanup := testutil.Logger(t)
	defer cleanup()

	opts := weshnet.TestingOpts{
		Mocknet:     libp2p_mocknet.New(),
		Logger:      logger,
		ConnectFunc: weshnet.ConnectAll,
	}

	nodes, cleanup := weshnet.NewTestingProtocolWithMockedPeers(ctx, t, &opts, nil, 2)
	defer cleanup()

	t.Log("addAsContact")
	addAsContact(ctx, t, nodes, nodes)

	// send messages before deactivating
	t.Log("sendMessageToContact")
	sendMessageToContact(ctx, t, []string{"pre-deactivate"}, nodes)

	// get contact group
	contactGroup := getContactGroup(ctx, t, nodes[0], nodes[1])

	// deactivate contact group
	deactivateFunc := func() {
		t.Log("DeactivateGroup")
		_, err := nodes[0].Client.DeactivateGroup(ctx, &protocoltypes.DeactivateGroup_Request{
			GroupPK: contactGroup.Group.PublicKey,
		})
		require.NoError(t, err)
	}

	// reactivate group
	activateFunc := func() {
		t.Log("ActivateGroup")
		_, err := nodes[0].Client.ActivateGroup(ctx, &protocoltypes.ActivateGroup_Request{
			GroupPK: contactGroup.Group.PublicKey,
		})
		require.NoError(t, err)
	}

	go deactivateFunc()
	time.Sleep(1 * time.Millisecond)
	go activateFunc()

	// send message after reactivating
	time.Sleep(5 * time.Second)
	t.Log("sendMessageToContact")
	sendMessageToContact(ctx, t, []string{"post-reactivate"}, nodes)
}

func TestReactivateMultimemberGroup(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	logger, cleanup := testutil.Logger(t)
	defer cleanup()

	opts := weshnet.TestingOpts{
		Mocknet:     libp2p_mocknet.New(),
		Logger:      logger,
		ConnectFunc: weshnet.ConnectAll,
	}

	nodes, cleanup := weshnet.NewTestingProtocolWithMockedPeers(ctx, t, &opts, nil, 2)
	defer cleanup()

	// Create MultiMember Group
	group := weshnet.CreateMultiMemberGroupInstance(ctx, t, nodes[0], nodes[1])

	// Send message before deactivation
	sendMessageOnGroup(ctx, t, nodes, nodes, group.PublicKey, []string{"pre-deactivate"})

	// deactivate multimember group
	_, err := nodes[0].Client.DeactivateGroup(ctx, &protocoltypes.DeactivateGroup_Request{
		GroupPK: group.PublicKey,
	})
	require.NoError(t, err)

	// reactivate group
	_, err = nodes[0].Client.ActivateGroup(ctx, &protocoltypes.ActivateGroup_Request{
		GroupPK: group.PublicKey,
	})
	require.NoError(t, err)

	// Send message after reactivation
	sendMessageOnGroup(ctx, t, nodes, nodes, group.PublicKey, []string{"post-deactivate"})
}
