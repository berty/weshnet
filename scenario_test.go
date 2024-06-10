//go:build !js

package weshnet_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	libp2p_mocknet "github.com/berty/go-libp2p-mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"go.uber.org/zap"

	weshnet "berty.tech/weshnet"
	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/testutil"
)

type testCase struct {
	Name           string
	NumberOfClient int
	ConnectFunc    weshnet.ConnectTestingProtocolFunc
	Speed          testutil.Speed
	Stability      testutil.Stability
	Timeout        time.Duration
}

type testFunc func(context.Context, *testing.T, ...*weshnet.TestingProtocol)

// Tests

func TestScenario_CreateMultiMemberGroup(t *testing.T) {
	cases := []testCase{
		{"2 clients/connectAll", 2, weshnet.ConnectAll, testutil.Fast, testutil.Flappy, time.Second * 10},
		{"3 clients/connectAll", 3, weshnet.ConnectAll, testutil.Fast, testutil.Flappy, time.Second * 10},
		{"3 clients/connectInLine", 3, weshnet.ConnectInLine, testutil.Fast, testutil.Flappy, time.Second * 10},
		{"5 clients/connectAll", 5, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 20},
		{"5 clients/connectInLine", 5, weshnet.ConnectInLine, testutil.Slow, testutil.Flappy, time.Second * 20},
		{"8 clients/connectAll", 8, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 60},
		{"8 clients/connectInLine", 8, weshnet.ConnectInLine, testutil.Slow, testutil.Flappy, time.Second * 60},
		{"10 clients/connectAll", 10, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 90},
		{"10 clients/connectInLine", 10, weshnet.ConnectInLine, testutil.Slow, testutil.Flappy, time.Second * 90},
	}

	testingScenario(t, cases, func(ctx context.Context, t *testing.T, tps ...*weshnet.TestingProtocol) {
		createMultiMemberGroup(ctx, t, tps...)
	})
}

func TestScenario_MessageMultiMemberGroup(t *testing.T) {
	cases := []testCase{
		{"2 clients/connectAll", 2, weshnet.ConnectAll, testutil.Fast, testutil.Flappy, time.Second * 10},
		{"3 clients/connectAll", 3, weshnet.ConnectAll, testutil.Fast, testutil.Flappy, time.Second * 10},
		{"3 clients/connectInLine", 3, weshnet.ConnectInLine, testutil.Fast, testutil.Flappy, time.Second * 10},
		{"5 clients/connectAll", 5, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 20},
		{"5 clients/connectInLine", 5, weshnet.ConnectInLine, testutil.Slow, testutil.Flappy, time.Second * 20},
		{"8 clients/connectAll", 8, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 40},
		{"8 clients/connectInLine", 8, weshnet.ConnectInLine, testutil.Slow, testutil.Flappy, time.Second * 40},
		{"10 clients/connectAll", 10, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 90},
		{"10 clients/connectInLine", 10, weshnet.ConnectInLine, testutil.Slow, testutil.Flappy, time.Second * 90},
	}

	testingScenario(t, cases, func(ctx context.Context, t *testing.T, tps ...*weshnet.TestingProtocol) {
		// Create MultiMember Group
		groupID := createMultiMemberGroup(ctx, t, tps...)

		// Each member sends 3 messages on MultiMember Group
		messages := []string{"test1", "test2", "test3"}
		sendMessageOnGroup(ctx, t, tps, tps, groupID, messages)
	})
}

func TestScenario_GroupDeviceStatusOnMultiMemberGroup(t *testing.T) {
	cases := []testCase{
		{"2 clients/connectAll", 2, weshnet.ConnectAll, testutil.Fast, testutil.Flappy, time.Second * 10},
		{"3 clients/connectAll", 3, weshnet.ConnectAll, testutil.Fast, testutil.Flappy, time.Second * 10},
		{"5 clients/connectAll", 5, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 20},
		{"8 clients/connectAll", 8, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 60},
		{"10 clients/connectAll", 10, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 90},
	}

	testingScenario(t, cases, func(ctx context.Context, t *testing.T, tps ...*weshnet.TestingProtocol) {
		// Create MultiMember Group
		groupID := createMultiMemberGroup(ctx, t, tps...)

		testGroupDeviceStatus(ctx, t, groupID, tps...)
	})
}

func testGroupDeviceStatus(ctx context.Context, t *testing.T, groupID []byte, tps ...*weshnet.TestingProtocol) {
	ntps := len(tps)

	// Get group device status
	{
		testutil.LogTree(t, "Get Group Device Status", 1, true)
		start := time.Now()

		wg := sync.WaitGroup{}
		statusReceivedLock := sync.Mutex{}
		statusReceived := make([]map[string]struct{}, ntps)
		wg.Add(ntps)

		nSuccess := int64(0)
		for i := range tps {
			go func(i int) {
				tp := tps[i]
				defer wg.Done()

				statusReceived[i] = map[string]struct{}{}

				ctx, cancel := context.WithCancel(ctx)
				defer cancel()

				sub, inErr := tp.Client.GroupDeviceStatus(ctx, &protocoltypes.GroupDeviceStatus_Request{
					GroupPK: groupID,
				})
				if inErr != nil {
					assert.NoError(t, inErr, fmt.Sprintf("error for client %d", i))
					return
				}

				for {
					evt, inErr := sub.Recv()
					if inErr != nil {
						if inErr != io.EOF {
							assert.NoError(t, inErr, fmt.Sprintf("error for client %d", i))
						}

						break
					}

					assert.Equal(t, evt.Type, protocoltypes.TypePeerConnected)
					connected := &protocoltypes.GroupDeviceStatus_Reply_PeerConnected{}
					err := connected.Unmarshal(evt.Event)
					assert.NoError(t, err, fmt.Sprintf("Unmarshal error for client %d", i))

					statusReceivedLock.Lock()
					statusReceived[i][connected.PeerID] = struct{}{}
					done := len(statusReceived[i]) == ntps-1
					statusReceivedLock.Unlock()

					if done {
						n := atomic.AddInt64(&nSuccess, 1)

						got := fmt.Sprintf("%d/%d", n, ntps)
						tps[i].Opts.Logger.Debug("received all group device status", zap.String("ok", got))
						return
					}
				}
			}(i)
		}

		wg.Wait()

		statusReceivedLock.Lock()
		ok := true
		for i := range statusReceived {
			if !assert.Equal(t, ntps-1, len(statusReceived[i]), fmt.Sprintf("mismatch for client %d", i)) {
				ok = false
			}
		}
		require.True(t, ok)
		statusReceivedLock.Unlock()

		testutil.LogTree(t, "duration: %s", 1, false, time.Since(start))
	}
}

//
//func TestScenario_MessageMultiMemberGroup2(t *testing.T) {
//	cases := []testCase{
//		{"2 clients/connectAll", 2, ConnectAll, testutil.Fast, testutil.Stable, time.Second * 60},
//	}
//
//	testingScenario(t, cases, func(ctx context.Context, t *testing.T, tps ...*TestingProtocol) {
//		// Create MultiMember Group
//		groupID := createMultiMemberGroup(ctx, t, tps...)
//
//		const messageCount = 100
//		// Each member sends 3 messages on MultiMember Group
//		messages := make([]string, messageCount)
//		for i := 0; i < messageCount; i++ {
//			messages[i] = fmt.Sprintf("test%d", i)
//		}
//
//		sendMessageOnGroup(ctx, t, tps, tps, groupID, messages)
//	})
//}

func TestScenario_MessageSeveralMultiMemberGroups(t *testing.T) {
	const ngroup = 3

	cases := []testCase{
		{"2 clients/connectAll", 2, weshnet.ConnectAll, testutil.Fast, testutil.Flappy, time.Second * 20},
		{"3 clients/connectAll", 3, weshnet.ConnectAll, testutil.Fast, testutil.Flappy, time.Second * 20},
		{"3 clients/connectInLine", 3, weshnet.ConnectInLine, testutil.Fast, testutil.Flappy, time.Second * 20},
		{"5 clients/connectAll", 5, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 60},
		{"5 clients/connectInLine", 5, weshnet.ConnectInLine, testutil.Slow, testutil.Flappy, time.Second * 60},
		{"8 clients/connectAll", 8, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 180},
		{"8 clients/connectInLine", 8, weshnet.ConnectInLine, testutil.Slow, testutil.Flappy, time.Second * 180},
		{"10 clients/connectAll", 10, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 300},
		{"10 clients/connectInLine", 10, weshnet.ConnectInLine, testutil.Slow, testutil.Flappy, time.Second * 300},
	}

	testingScenario(t, cases, func(ctx context.Context, t *testing.T, tps ...*weshnet.TestingProtocol) {
		for i := 0; i < ngroup; i++ {
			t.Logf("===== MultiMember Group #%d =====", i+1)
			// Create MultiMember Group
			groupID := createMultiMemberGroup(ctx, t, tps...)

			// Each member sends 3 messages on MultiMember Group
			messages := []string{"test1", "test2", "test3"}
			sendMessageOnGroup(ctx, t, tps, tps, groupID, messages)
		}
	})
}

func TestScenario_AddContact(t *testing.T) {
	cases := []testCase{
		{"2 clients/connectAll", 2, weshnet.ConnectAll, testutil.Fast, testutil.Flappy, time.Second * 20},
		{"3 clients/connectAll", 3, weshnet.ConnectAll, testutil.Fast, testutil.Flappy, time.Second * 20},
		{"5 clients/connectAll", 5, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 30},
		{"8 clients/connectAll", 8, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 40},
		{"10 clients/connectAll", 10, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 60},
	}

	testingScenario(t, cases, func(ctx context.Context, t *testing.T, tps ...*weshnet.TestingProtocol) {
		addAsContact(ctx, t, tps, tps)
	})
}

func TestScenario_MessageContactGroup(t *testing.T) {
	cases := []testCase{
		{"2 clients/connectAll", 2, weshnet.ConnectAll, testutil.Fast, testutil.Flappy, time.Second * 20},
		{"3 clients/connectAll", 3, weshnet.ConnectAll, testutil.Fast, testutil.Flappy, time.Second * 20},
		{"5 clients/connectAll", 5, weshnet.ConnectAll, testutil.Slow, testutil.Flappy, time.Second * 30},
		{"8 clients/connectAll", 8, weshnet.ConnectAll, testutil.Slow, testutil.Broken, time.Second * 40},
		{"10 clients/connectAll", 10, weshnet.ConnectAll, testutil.Slow, testutil.Broken, time.Second * 60},
	}

	testingScenario(t, cases, func(ctx context.Context, t *testing.T, tps ...*weshnet.TestingProtocol) {
		// Add accounts as contacts
		addAsContact(ctx, t, tps, tps)

		// Send messages between all accounts on contact groups
		messages := []string{"test1", "test2", "test3"}
		sendMessageToContact(ctx, t, messages, tps)
	})
}

func TestScenario_MessageAccountGroup(t *testing.T) {
	cases := []testCase{
		{"1 client/connectAll", 1, weshnet.ConnectAll, testutil.Fast, testutil.Stable, time.Second * 10},
	}

	testingScenario(t, cases, func(ctx context.Context, t *testing.T, tps ...*weshnet.TestingProtocol) {
		// Get account config
		config, err := tps[0].Client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
		require.NoError(t, err)
		require.NotNil(t, config)

		// Send messages on account group
		messages := []string{"test1", "test2", "test3"}
		sendMessageOnGroup(ctx, t, tps, tps, config.AccountGroupPK, messages)
	})
}

func TestScenario_MessageAccountGroup_NonMocked(t *testing.T) {
	cases := []testCase{
		{"1 client/connectAll", 1, weshnet.ConnectAll, testutil.Fast, testutil.Stable, time.Second * 10},
	}

	testingScenarioNonMocked(t, cases, func(ctx context.Context, t *testing.T, tps ...*weshnet.TestingProtocol) {
		// Get account config
		config, err := tps[0].Client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
		require.NoError(t, err)
		require.NotNil(t, config)

		// Send messages on account group
		messages := []string{"test1", "test2", "test3"}
		sendMessageOnGroup(ctx, t, tps, tps, config.AccountGroupPK, messages)
	})
}

func TestScenario_MessageAccountAndMultiMemberGroups(t *testing.T) {
	cases := []testCase{
		{"2 clients/connectAll", 2, weshnet.ConnectAll, testutil.Fast, testutil.Broken, time.Second * 10},
		{"3 clients/connectAll", 3, weshnet.ConnectAll, testutil.Fast, testutil.Broken, time.Second * 10},
		{"3 clients/connectInLine", 3, weshnet.ConnectInLine, testutil.Fast, testutil.Broken, time.Second * 10},
		{"5 clients/connectAll", 5, weshnet.ConnectAll, testutil.Slow, testutil.Broken, time.Second * 20},
		{"5 clients/connectInLine", 5, weshnet.ConnectInLine, testutil.Slow, testutil.Broken, time.Second * 20},
		{"8 clients/connectAll", 8, weshnet.ConnectAll, testutil.Slow, testutil.Broken, time.Second * 30},
		{"8 clients/connectInLine", 8, weshnet.ConnectInLine, testutil.Slow, testutil.Broken, time.Second * 30},
		{"10 clients/connectAll", 10, weshnet.ConnectAll, testutil.Slow, testutil.Broken, time.Second * 40},
		{"10 clients/connectInLine", 10, weshnet.ConnectInLine, testutil.Slow, testutil.Broken, time.Second * 40},
	}

	testingScenario(t, cases, func(ctx context.Context, t *testing.T, tps ...*weshnet.TestingProtocol) {
		t.Log("===== Send Messages on MultiMember Group =====")
		// Create MultiMember Group
		mmGroup := createMultiMemberGroup(ctx, t, tps...)

		// Each member sends 3 messages on MultiMember Group
		messages := []string{"test1", "test2", "test3"}
		sendMessageOnGroup(ctx, t, tps, tps, mmGroup, messages)

		t.Log("===== Send Messages on Account Group =====")
		// Send messages on account groups
		for _, account := range tps {
			// Get account config
			config, err := account.Client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
			require.NoError(t, err)
			require.NotNil(t, config)

			// Send messages on account group
			messages = []string{"account1", "account2", "account3"}
			sendMessageOnGroup(ctx, t, []*weshnet.TestingProtocol{account}, []*weshnet.TestingProtocol{account}, config.AccountGroupPK, messages)
		}

		t.Log("===== Send Messages again on MultiMember Group =====")
		// Each member sends 3 messages on MultiMember Group
		messages = []string{"test4", "test5", "test6"}
		sendMessageOnGroup(ctx, t, tps, tps, mmGroup, messages)
	})
}

func TestScenario_MessageAccountAndContactGroups(t *testing.T) {
	cases := []testCase{
		{"2 clients/connectAll", 2, weshnet.ConnectAll, testutil.Fast, testutil.Broken, time.Second * 10},
		{"3 clients/connectAll", 3, weshnet.ConnectAll, testutil.Fast, testutil.Broken, time.Second * 10},
		{"5 clients/connectAll", 5, weshnet.ConnectAll, testutil.Slow, testutil.Broken, time.Second * 20},
		{"8 clients/connectAll", 8, weshnet.ConnectAll, testutil.Slow, testutil.Broken, time.Second * 30},
		{"10 clients/connectAll", 10, weshnet.ConnectAll, testutil.Slow, testutil.Broken, time.Second * 40},
	}

	testingScenario(t, cases, func(ctx context.Context, t *testing.T, tps ...*weshnet.TestingProtocol) {
		t.Log("===== Send Messages on Contact Group =====")
		// Add accounts as contacts
		addAsContact(ctx, t, tps, tps)
		// Send messages between all accounts on contact groups
		messages := []string{"contact1", "contact2", "contact3"}
		sendMessageToContact(ctx, t, messages, tps)

		t.Log("===== Send Messages on Account Group =====")
		// Send messages on account groups
		for _, account := range tps {
			// Get account config
			config, err := account.Client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
			require.NoError(t, err)
			require.NotNil(t, config)

			// Send messages on account group
			messages = []string{"account1", "account2", "account3"}
			sendMessageOnGroup(ctx, t, []*weshnet.TestingProtocol{account}, []*weshnet.TestingProtocol{account}, config.AccountGroupPK, messages)
		}

		t.Log("===== Send Messages again on Contact Group =====")
		// Send messages between all accounts on contact groups
		messages = []string{"contact4", "contact5", "contact6"}
		sendMessageToContact(ctx, t, messages, tps)
	})
}

// Helpers

func testingScenario(t *testing.T, tcs []testCase, tf testFunc) {
	if os.Getenv("WITH_GOLEAK") == "1" {
		defer goleak.VerifyNone(t,
			goleak.IgnoreTopFunction("github.com/syndtr/goleveldb/leveldb.(*DB).mpoolDrain"),     // inherited from one of the imports (init)
			goleak.IgnoreTopFunction("github.com/ipfs/go-log/writer.(*MirrorWriter).logRoutine"), // inherited from one of the imports (init)
			goleak.IgnoreTopFunction("github.com/jbenet/goprocess/periodic.callOnTicker.func1"),  // inherited from github.com/ipfs/kubo/core.NewNode
			goleak.IgnoreTopFunction("go.opencensus.io/stats/view.(*worker).start"),              // inherited from github.com/ipfs/kubo/core.NewNode)
			goleak.IgnoreTopFunction("github.com/desertbit/timer.timerRoutine"),                  // inherited from github.com/ipfs/kubo/core.NewNode)
			goleak.IgnoreTopFunction("go.opentelemetry.io/otel/instrumentation/grpctrace.wrapClientStream.func1"),
			goleak.IgnoreTopFunction("go.opentelemetry.io/otel/instrumentation/grpctrace.StreamClientInterceptor.func1.1"),
		)
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			testutil.FilterStabilityAndSpeed(t, tc.Stability, tc.Speed)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			logger, cleanup := testutil.Logger(t)
			defer cleanup()

			mn := libp2p_mocknet.New()
			defer mn.Close()

			opts := weshnet.TestingOpts{
				Mocknet:     mn,
				Logger:      logger,
				ConnectFunc: tc.ConnectFunc,
			}

			tps, cleanup := weshnet.NewTestingProtocolWithMockedPeers(ctx, t, &opts, nil, tc.NumberOfClient)
			defer cleanup()

			var cctx context.Context

			if tc.Timeout > 0 {
				cctx, cancel = context.WithTimeout(ctx, tc.Timeout)
			} else {
				cctx, cancel = context.WithCancel(ctx)
			}

			tf(cctx, t, tps...)
			cancel()
		})
	}
}

func testingScenarioNonMocked(t *testing.T, tcs []testCase, tf testFunc) {
	if os.Getenv("WITH_GOLEAK") == "1" {
		defer goleak.VerifyNone(t,
			goleak.IgnoreTopFunction("github.com/syndtr/goleveldb/leveldb.(*DB).mpoolDrain"),     // inherited from one of the imports (init)
			goleak.IgnoreTopFunction("github.com/ipfs/go-log/writer.(*MirrorWriter).logRoutine"), // inherited from one of the imports (init)
			goleak.IgnoreTopFunction("github.com/jbenet/goprocess/periodic.callOnTicker.func1"),  // inherited from github.com/ipfs/kubo/core.NewNode
			goleak.IgnoreTopFunction("go.opencensus.io/stats/view.(*worker).start"),              // inherited from github.com/ipfs/kubo/core.NewNode)
			goleak.IgnoreTopFunction("github.com/desertbit/timer.timerRoutine"),                  // inherited from github.com/ipfs/kubo/core.NewNode)
			goleak.IgnoreTopFunction("go.opentelemetry.io/otel/instrumentation/grpctrace.wrapClientStream.func1"),
			goleak.IgnoreTopFunction("go.opentelemetry.io/otel/instrumentation/grpctrace.StreamClientInterceptor.func1.1"),
		)
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			testutil.FilterStabilityAndSpeed(t, tc.Stability, tc.Speed)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			logger, cleanup := testutil.Logger(t)
			defer cleanup()

			mn := libp2p_mocknet.New()
			defer mn.Close()

			opts := weshnet.TestingOpts{
				Mocknet:     mn,
				Logger:      logger,
				ConnectFunc: tc.ConnectFunc,
			}

			tps, cleanup := weshnet.NewTestingProtocolWithMockedPeers(ctx, t, &opts, nil, tc.NumberOfClient)
			defer cleanup()

			var cctx context.Context

			if tc.Timeout > 0 {
				cctx, cancel = context.WithTimeout(ctx, tc.Timeout)
			} else {
				cctx, cancel = context.WithCancel(ctx)
			}

			tf(cctx, t, tps...)
			cancel()
		})
	}
}

func createMultiMemberGroup(ctx context.Context, t *testing.T, tps ...*weshnet.TestingProtocol) (groupID []byte) {
	return weshnet.CreateMultiMemberGroupInstance(ctx, t, tps...).PublicKey
}

func addAsContact(ctx context.Context, t *testing.T, senders, receivers []*weshnet.TestingProtocol) {
	testutil.LogTree(t, "Add Senders/Receivers as Contact", 0, true)
	start := time.Now()
	var sendDuration, receiveDuration, acceptDuration, activateDuration time.Duration

	for i, sender := range senders {
		for _, receiver := range receivers {
			substart := time.Now()

			// Get sender/receiver configs
			senderCfg, err := sender.Client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
			require.NoError(t, err)
			require.NotNil(t, senderCfg)
			receiverCfg, err := receiver.Client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
			require.NoError(t, err)
			require.NotNil(t, receiverCfg)

			// Setup receiver's shareable contact
			var receiverRDVSeed []byte

			crf, err := receiver.Client.ContactRequestReference(ctx, &protocoltypes.ContactRequestReference_Request{})
			if err != nil || !crf.Enabled || len(crf.PublicRendezvousSeed) == 0 {
				_, err = receiver.Client.ContactRequestEnable(ctx, &protocoltypes.ContactRequestEnable_Request{})
				require.NoError(t, err)
				receiverRDV, err := receiver.Client.ContactRequestResetReference(ctx, &protocoltypes.ContactRequestResetReference_Request{})
				require.NoError(t, err)
				require.NotNil(t, receiverRDV)
				receiverRDVSeed = receiverRDV.PublicRendezvousSeed
			} else {
				receiverRDVSeed = crf.PublicRendezvousSeed
			}

			receiverSharableContact := &protocoltypes.ShareableContact{
				PK:                   receiverCfg.AccountPK,
				PublicRendezvousSeed: receiverRDVSeed,
			}

			// Sender sends contact request
			_, err = sender.Client.ContactRequestSend(ctx, &protocoltypes.ContactRequestSend_Request{
				Contact: receiverSharableContact,
			})

			// Check if sender and receiver are the same account, should return the right error and skip
			if bytes.Equal(senderCfg.AccountPK, receiverCfg.AccountPK) {
				require.Equal(t, errcode.LastCode(err), errcode.ErrContactRequestSameAccount)
				continue
			}

			// Check if contact request was already sent, should return right error and skip
			receiverWasSender := false
			for j := 0; j < i; j++ {
				if senders[j] == receiver {
					receiverWasSender = true
				}
			}

			senderWasReceiver := false
			if receiverWasSender {
				for _, r := range receivers {
					if r == sender {
						senderWasReceiver = true
					}
				}
			}

			if receiverWasSender && senderWasReceiver {
				require.Equal(t, errcode.LastCode(err), errcode.ErrContactRequestContactAlreadyAdded)
				continue
			}

			// No other error should occur
			require.NoError(t, err)

			sendDuration += time.Since(substart)
			substart = time.Now()

			// Receiver subscribes to handle incoming contact request
			subCtx, subCancel := context.WithCancel(ctx)
			subReceiver, err := receiver.Client.GroupMetadataList(subCtx, &protocoltypes.GroupMetadataList_Request{
				GroupPK: receiverCfg.AccountGroupPK,
			})
			require.NoError(t, err)
			found := false

			// Receiver waits for valid contact request coming from sender
			for {
				evt, err := subReceiver.Recv()
				if err == io.EOF || subReceiver.Context().Err() != nil {
					break
				}

				require.NoError(t, err)

				if evt == nil || evt.Metadata.EventType != protocoltypes.EventTypeAccountContactRequestIncomingReceived {
					continue
				}

				req := &protocoltypes.AccountContactRequestIncomingReceived{}
				err = req.Unmarshal(evt.Event)

				require.NoError(t, err)

				if bytes.Equal(senderCfg.AccountPK, req.ContactPK) {
					found = true
					break
				}
			}

			subCancel()
			require.True(t, found)

			receiveDuration += time.Since(substart)
			substart = time.Now()

			// Receiver accepts contact request
			_, err = receiver.Client.ContactRequestAccept(ctx, &protocoltypes.ContactRequestAccept_Request{
				ContactPK: senderCfg.AccountPK,
			})

			require.NoError(t, err)

			acceptDuration += time.Since(substart)
			substart = time.Now()

			// Both receiver and sender activate the contact group
			grpInfo, err := sender.Client.GroupInfo(ctx, &protocoltypes.GroupInfo_Request{
				ContactPK: receiverCfg.AccountPK,
			})
			require.NoError(t, err)

			_, err = sender.Client.ActivateGroup(ctx, &protocoltypes.ActivateGroup_Request{
				GroupPK: grpInfo.Group.PublicKey,
			})

			require.NoError(t, err)

			grpInfo2, err := receiver.Client.GroupInfo(ctx, &protocoltypes.GroupInfo_Request{
				ContactPK: senderCfg.AccountPK,
			})
			require.NoError(t, err)

			require.Equal(t, grpInfo.Group.PublicKey, grpInfo2.Group.PublicKey)

			_, err = receiver.Client.ActivateGroup(ctx, &protocoltypes.ActivateGroup_Request{
				GroupPK: grpInfo2.Group.PublicKey,
			})

			require.NoError(t, err)

			activateDuration += time.Since(substart)
		}
	}

	testutil.LogTree(t, "Send Contact Requests", 1, true)
	testutil.LogTree(t, "duration: %s", 1, false, sendDuration)
	testutil.LogTree(t, "Receive Contact Requests", 1, true)
	testutil.LogTree(t, "duration: %s", 1, false, receiveDuration)
	testutil.LogTree(t, "Accept Contact Requests", 1, true)
	testutil.LogTree(t, "duration: %s", 1, false, acceptDuration)
	testutil.LogTree(t, "Activate Contact Groups", 1, true)
	testutil.LogTree(t, "duration: %s", 1, false, activateDuration)

	testutil.LogTree(t, "duration: %s", 0, false, time.Since(start))
}

func getContactGroup(ctx context.Context, t *testing.T, source *weshnet.TestingProtocol, contact *weshnet.TestingProtocol) *protocoltypes.GroupInfo_Reply {
	// Get contact group
	contactGroup, err := source.Client.GroupInfo(ctx, &protocoltypes.GroupInfo_Request{
		ContactPK: getAccountPubKey(t, contact),
	})
	require.NoError(t, err)
	require.NotNil(t, contactGroup)
	return contactGroup
}

func sendMessageToContact(ctx context.Context, t *testing.T, messages []string, tps []*weshnet.TestingProtocol) {
	for _, sender := range tps {
		for _, receiver := range tps {
			// Don't try to send messages to itself using contact group
			if sender == receiver {
				continue
			}

			// Get contact group
			contactGroup := getContactGroup(ctx, t, sender, receiver)

			// Send messages on contact group
			sendMessageOnGroup(ctx, t, []*weshnet.TestingProtocol{sender}, []*weshnet.TestingProtocol{receiver}, contactGroup.Group.PublicKey, messages)
		}
	}
}

func sendMessageOnGroup(ctx context.Context, t *testing.T, senders, receivers []*weshnet.TestingProtocol, groupPK []byte, messages []string) {
	testutil.LogTree(t, "Send, Receive and List Messages", 0, true)
	start := time.Now()

	// Setup expectedMessages map
	expectedMessages := map[string]struct{}{}
	expectedMessagesCount := len(messages) * len(senders)
	expectedMessagesLock := sync.Mutex{}

	for _, message := range messages {
		for _, sender := range senders {
			expectedMessage := getAccountB64PubKey(t, sender) + " - " + message
			expectedMessages[expectedMessage] = struct{}{}
		}
	}

	// Setup map to check expected messages reception
	subReceivedMessages := map[string]map[string]bool{}
	subReceivedMessagesCount := map[string]int{}
	listReceivedMessages := map[string]map[string]bool{}
	listReceivedMessagesCount := map[string]int{}

	for _, receiver := range receivers {
		subReceiverMap := map[string]bool{}
		listReceiverMap := map[string]bool{}

		for expectedMessage := range expectedMessages {
			subReceiverMap[expectedMessage] = false
			listReceiverMap[expectedMessage] = false
		}

		receiverID := getAccountB64PubKey(t, receiver)
		subReceivedMessages[receiverID] = subReceiverMap
		listReceivedMessages[receiverID] = listReceiverMap
		subReceivedMessagesCount[receiverID] = 0
		listReceivedMessagesCount[receiverID] = 0
	}
	receivedMessagesLock := sync.Mutex{}

	// Senders send all expected messages
	{
		testutil.LogTree(t, "Senders Send Messages", 1, true)
		start := time.Now()

		for _, sender := range senders {
			senderID := getAccountB64PubKey(t, sender)
			for _, message := range messages {
				_, err := sender.Client.AppMessageSend(ctx, &protocoltypes.AppMessageSend_Request{
					GroupPK: groupPK,
					Payload: []byte(senderID + " - " + message),
				})

				require.NoError(t, err)
			}
		}

		testutil.LogTree(t, "duration: %s", 1, false, time.Since(start))
	}

	// Receivers receive all expected messages
	{
		testutil.LogTree(t, "Receivers Receive Messages (subscription)", 1, true)
		start := time.Now()

		var wg sync.WaitGroup
		wg.Add(len(receivers))

		for _, receiver := range receivers {
			// Subscribe receivers to wait for incoming messages
			go func(receiver *weshnet.TestingProtocol) {
				subCtx, subCancel := context.WithCancel(ctx)
				defer subCancel()
				defer wg.Done()

				sub, err := receiver.Client.GroupMessageList(subCtx, &protocoltypes.GroupMessageList_Request{
					GroupPK: groupPK,
				})
				if !assert.NoError(t, err) {
					return
				}

				receiverID := getAccountB64PubKey(t, receiver)

				for {
					if subCtx.Err() != nil {
						return
					}

					// Receive message
					res, err := sub.Recv()
					if err == io.EOF {
						return
					}
					if !assert.NoError(t, err) {
						continue
					}

					// Check if received message was expected
					expectedMessagesLock.Lock()
					_, expected := expectedMessages[string(res.Message)]
					expectedMessagesLock.Unlock()
					if !expected {
						continue
					}

					// Check if message was already received
					receivedMessagesLock.Lock()
					alreadyReceived := subReceivedMessages[receiverID][string(res.Message)]
					if alreadyReceived {
						receivedMessagesLock.Unlock()
						continue
					}

					// Mark message as received
					subReceivedMessages[receiverID][string(res.Message)] = true
					subReceivedMessagesCount[receiverID]++
					// Return if all expected messages were received
					if subReceivedMessagesCount[receiverID] == expectedMessagesCount {
						receivedMessagesLock.Unlock()
						return
					}
					receivedMessagesLock.Unlock()
				}
			}(receiver)
		}

		// Wait that all receivers received messages
		wg.Wait()

		// Check if everything is ok
		for _, receiver := range receivers {
			receiverID := getAccountB64PubKey(t, receiver)
			assert.Equal(t, expectedMessagesCount, subReceivedMessagesCount[receiverID])
		}

		testutil.LogTree(t, "duration: %s", 1, false, time.Since(start))
	}

	// Receivers list all expected messages
	{
		testutil.LogTree(t, "Receivers List Messages (store)", 1, true)
		start := time.Now()

		var wg sync.WaitGroup
		wg.Add(len(receivers))

		for _, receiver := range receivers {
			// Subscribe receivers to wait for incoming messages
			go func(receiver *weshnet.TestingProtocol) {
				subCtx, subCancel := context.WithCancel(ctx)
				defer subCancel()
				defer wg.Done()

				req := protocoltypes.GroupMessageList_Request{
					GroupPK:  groupPK,
					UntilNow: true,
				}

				ml, err := receiver.Client.GroupMessageList(subCtx, &req)
				if !assert.NoError(t, err) {
					return
				}

				receiverID := getAccountB64PubKey(t, receiver)

				for {
					if subCtx.Err() != nil {
						return
					}

					// Receive message
					res, err := ml.Recv()
					if err == io.EOF {
						return
					}
					if !assert.NoError(t, err) {
						continue
					}

					// Check if received message was expected
					expectedMessagesLock.Lock()
					_, expected := expectedMessages[string(res.Message)]
					expectedMessagesLock.Unlock()
					if !expected {
						continue
					}

					// Check if message was already received
					receivedMessagesLock.Lock()
					alreadyReceived := listReceivedMessages[receiverID][string(res.Message)]
					if alreadyReceived {
						receivedMessagesLock.Unlock()
						continue
					}

					// Mark message as received
					listReceivedMessages[receiverID][string(res.Message)] = true
					listReceivedMessagesCount[receiverID]++
					// Return if all expected messages were received
					if listReceivedMessagesCount[receiverID] == expectedMessagesCount {
						receivedMessagesLock.Unlock()
						return
					}
					receivedMessagesLock.Unlock()
				}
			}(receiver)
		}

		// Wait that all receivers listed messages
		wg.Wait()

		// Check if everything is ok
		for _, receiver := range receivers {
			receiverID := getAccountB64PubKey(t, receiver)
			assert.Equal(t, expectedMessagesCount, listReceivedMessagesCount[receiverID])
		}

		testutil.LogTree(t, "duration: %s", 1, false, time.Since(start))
	}
	testutil.LogTree(t, "duration: %s", 0, false, time.Since(start))
}

func getAccountPubKey(t *testing.T, tp *weshnet.TestingProtocol) []byte {
	t.Helper()

	_, accMemberDevice, err := tp.Opts.SecretStore.GetGroupForAccount()
	require.NoError(t, err)

	publicKeyBytes, err := accMemberDevice.Member().Raw()
	require.NoError(t, err)

	return publicKeyBytes
}

func getAccountB64PubKey(t *testing.T, tp *weshnet.TestingProtocol) string {
	t.Helper()

	tpPK := getAccountPubKey(t, tp)

	return base64.StdEncoding.EncodeToString(tpPK)
}
