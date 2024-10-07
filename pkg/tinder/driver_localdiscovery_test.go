package tinder

import (
	"context"
	"math/rand"
	"testing"
	"time"

	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/stretchr/testify/require"

	"berty.tech/weshnet/v2/pkg/testutil"
)

func TestServiceLocalDiscorvery(t *testing.T) {
	ctx := context.Background()
	mn := mocknet.New()
	defer mn.Close()

	logger, cleanup := testutil.Logger(t)
	defer cleanup()

	p1 := genLocalPeer(t, mn)
	p2 := genLocalPeer(t, mn)

	err := mn.LinkAll()
	require.NoError(t, err)

	disc1, err := NewLocalDiscovery(logger, p1, rand.New(rand.NewSource(rand.Int63())))
	require.NoError(t, err)

	disc2, err := NewLocalDiscovery(logger, p2, rand.New(rand.NewSource(rand.Int63())))
	require.NoError(t, err)

	s1, err := NewService(p1, logger, disc1)
	require.NoError(t, err)

	s2, err := NewService(p2, logger, disc2)
	require.NoError(t, err)

	const topic = "test_topic"

	s1.StartAdvertises(ctx, topic)

	// try a first lookup, should find nothing
	{
		out := s2.FindPeers(ctx, topic)
		peers := testPeersChanToSlice(t, out)
		require.Len(t, peers, 0, "no peer should be available")
	}

	{
		// start a subscribe and wait for connection
		sub := s2.Subscribe(topic)
		defer sub.Close()

		err = mn.ConnectAllButSelf()
		require.NoError(t, err)

		p, err := testWaitForPeers(t, sub.Out(), time.Second*5)
		require.NoError(t, err)
		require.Equal(t, p1.ID(), p.ID)
		require.Equal(t, p1.Addrs(), p.Addrs)
	}

	// try a lookup again, this time we should have some peers
	{
		out := s2.FindPeers(ctx, topic)
		p, err := testWaitForPeers(t, out, time.Second*5)

		require.NoError(t, err)
		require.Equal(t, p1.ID(), p.ID)
		require.Equal(t, p1.Addrs(), p.Addrs)
	}
}

func TestServiceLocalDiscorveryBeforeProtocolRegister(t *testing.T) {
	const topic = "test_topic"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mn := mocknet.New()
	defer mn.Close()

	logger, cleanup := testutil.Logger(t)
	defer cleanup()

	p1 := genLocalPeer(t, mn)
	p2 := genLocalPeer(t, mn)

	err := mn.LinkAll()
	require.NoError(t, err)

	disc1, err := NewLocalDiscovery(logger, p1, rand.New(rand.NewSource(rand.Int63())))
	require.NoError(t, err)

	// create service for peer 1
	s1, err := NewService(p1, logger, disc1)
	require.NoError(t, err)

	// start advertising
	s1.StartAdvertises(ctx, topic)

	// connect both peer BEFORE registering local discovery protocol for peer 2
	err = mn.ConnectAllButSelf()
	require.NoError(t, err)

	// let some time to peers to connect and trigger protocol exchange
	time.Sleep(time.Millisecond * 200)

	// register p2 local discovery
	disc2, err := NewLocalDiscovery(logger, p2, rand.New(rand.NewSource(rand.Int63())))
	require.NoError(t, err)

	s2, err := NewService(p2, logger, disc2)
	require.NoError(t, err)

	// start subscribe and wait for connection
	sub := s2.Subscribe(topic)
	defer sub.Close()

	// pull to fetch current peers on the topic
	sub.Pull()

	select {
	case p := <-sub.Out():
		require.Equal(t, p.ID, p1.ID())
	case <-time.After(time.Second * 5):
		require.FailNow(t, "unable to wait for peer on local discovery")
	}
}
