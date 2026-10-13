package weshnet_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/stretchr/testify/require"

	weshnet "berty.tech/weshnet/v2"
	"berty.tech/weshnet/v2/pkg/protocoltypes"
	"berty.tech/weshnet/v2/pkg/testutil"
)

// TestReplicationConvergence_OfflinePeerCatchesUp asserts that a peer which
// goes offline, misses messages and then reconnects catches up on every missed
// message in a single orbit-db head exchange — with no application-level retry.
//
// It runs the full weshnet protocol service (NewTestingProtocolWithMockedPeers).
// Note the one mocknet-imposed compromise: mocknet does not run the discovery
// poll + BackoffConnector that re-dials a dropped peer on a real network (a
// probe confirms a peer never reconnects on its own after the link is
// restored), so the test re-establishes the connection itself with a single
// ConnectPeers — modelling the one network reconnect that discovery would
// produce. Everything after that is the real stack: the reconnect re-fires the
// head exchange, which backfills the whole log. No retry loop is used.
//
// A single reconnect is reliable here because this isolated, clean reconnect
// rarely trips the go-libp2p-pubsub connect/disconnect race that can drop a
// peer's subscriptions on reconnect (the reconnect's new-peer event is deduped
// against the not-yet-removed stale peer entry). It is low-probability, not
// immune: the busier bertyreplication topology hits it ~1/100 and needs a
// reconnect retry loop.
func TestScenario_ReplicationConvergenceOfflinePeerCatchesUp(t *testing.T) {
	testutil.FilterStabilityAndSpeed(t, testutil.Flappy, testutil.Slow)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, cleanupLog := testutil.Logger(t)
	defer cleanupLog()

	mn := mocknet.New()
	defer mn.Close()

	opts := weshnet.TestingOpts{
		Mocknet:     mn,
		Logger:      logger,
		ConnectFunc: weshnet.ConnectAll,
	}
	tps, cleanupTps := weshnet.NewTestingProtocolWithMockedPeers(ctx, t, &opts, nil, 2)
	defer cleanupTps()

	alice, bob := tps[0], tps[1]
	aliceID, bobID := alice.Opts.Host.ID(), bob.Opts.Host.ID()

	// Create, join and activate the group on both peers.
	group := weshnet.CreateMultiMemberGroupInstance(ctx, t, alice, bob)
	groupPK := group.PublicKey

	// Baseline: bob receives a message from alice while connected, proving the
	// group is live before we partition them.
	require.NoError(t, sendMessages(ctx, alice, groupPK, []string{"online-1"}))
	waitForMessages(ctx, t, bob, groupPK, []string{"online-1"}, time.Second*20)

	// Take bob offline: sever the mocknet link (so no re-dial can succeed) and
	// close the existing connection.
	require.NoError(t, mn.UnlinkPeers(aliceID, bobID))
	require.NoError(t, mn.DisconnectPeers(aliceID, bobID))

	// Alice keeps sending while bob is offline.
	const offlineCount = 20
	offline := make([]string, offlineCount)
	for i := range offline {
		offline[i] = fmt.Sprintf("offline-%d", i)
	}
	require.NoError(t, sendMessages(ctx, alice, groupPK, offline))

	// Bring bob back online: restore the link and re-establish the connection
	// once. mocknet won't re-dial on its own, so this single ConnectPeers stands
	// in for the discovery-driven reconnect; we do not retry it, and we never
	// touch orbit-db. The reconnect re-fires the head exchange on its own.
	_, err := mn.LinkPeers(aliceID, bobID)
	require.NoError(t, err)
	_, err = mn.ConnectPeers(aliceID, bobID)
	require.NoError(t, err)

	// Assert bob catches up on every message it missed, with no manual help
	// beyond that single reconnect and no head-exchange retry.
	waitForMessages(ctx, t, bob, groupPK, offline, time.Second*60)
}

func sendMessages(ctx context.Context, tp *weshnet.TestingProtocol, groupPK []byte, messages []string) error {
	for _, m := range messages {
		if _, err := tp.Client.AppMessageSend(ctx, &protocoltypes.AppMessageSend_Request{
			GroupPk: groupPK,
			Payload: []byte(m),
		}); err != nil {
			return err
		}
	}
	return nil
}

// waitForMessages subscribes to the group and blocks until every wanted payload
// has been observed or the timeout fires. GroupMessageList replays the store
// and then tails new events, so messages backfilled after the subscription
// starts are streamed too.
func waitForMessages(ctx context.Context, t *testing.T, tp *weshnet.TestingProtocol, groupPK []byte, want []string, timeout time.Duration) {
	t.Helper()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	remaining := map[string]struct{}{}
	for _, m := range want {
		remaining[m] = struct{}{}
	}

	sub, err := tp.Client.GroupMessageList(ctx, &protocoltypes.GroupMessageList_Request{
		GroupPk: groupPK,
	})
	require.NoError(t, err)

	for len(remaining) > 0 {
		res, err := sub.Recv()
		if err != nil {
			require.NoError(t, err, "peer did not catch up, still missing %d/%d message(s)", len(remaining), len(want))
			return
		}
		delete(remaining, string(res.Message))
	}
}
