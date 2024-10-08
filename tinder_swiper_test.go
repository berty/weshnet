package weshnet

import (
	"context"
	"fmt"
	"testing"
	"time"

	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"berty.tech/weshnet/v2/pkg/ipfsutil"
	"berty.tech/weshnet/v2/pkg/rendezvous"
	"berty.tech/weshnet/v2/pkg/testutil"
	"berty.tech/weshnet/v2/pkg/tinder"
)

func TestAnnounceWatchForPeriod(t *testing.T) {
	testutil.FilterSpeed(t, testutil.Slow)
	cases := []struct {
		expectedPeersFound int
		topicA             []byte
		topicB             []byte
		seedA              []byte
		seedB              []byte
	}{
		{
			expectedPeersFound: 0,
			topicA:             []byte("topicA"),
			topicB:             []byte("topicB"),
			seedA:              []byte("seedA"),
			seedB:              []byte("seedA"),
		},
		{
			expectedPeersFound: 1,
			topicA:             []byte("topicA"),
			topicB:             []byte("topicA"),
			seedA:              []byte("seedA"),
			seedB:              []byte("seedA"),
		},
	}

	logger, cleanup := testutil.Logger(t)
	defer cleanup()
	for i, tc := range cases {
		t.Run(fmt.Sprintf("tc: %d", i), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			mn := mocknet.New()
			defer mn.Close()

			opts := &ipfsutil.TestingAPIOpts{
				Logger:          logger,
				Mocknet:         mn,
				DiscoveryServer: tinder.NewMockDriverServer(),
			}

			apiA := ipfsutil.TestingCoreAPIUsingMockNet(ctx, t, opts)
			apiB := ipfsutil.TestingCoreAPIUsingMockNet(ctx, t, opts)

			err := mn.LinkAll()
			require.NoError(t, err)
			err = mn.ConnectAllButSelf()
			require.NoError(t, err)

			rpA := rendezvous.NewRotationInterval(time.Hour)
			rpB := rendezvous.NewRotationInterval(time.Hour)

			swiperA := NewSwiper(opts.Logger, apiA.Tinder(), rpA)
			swiperB := NewSwiper(opts.Logger, apiB.Tinder(), rpB)

			swiperA.Announce(ctx, tc.topicA, tc.seedA)

			time.Sleep(time.Millisecond * 100)

			cpeers := swiperB.WatchTopic(ctx, tc.topicB, tc.seedB)

			var foundPeers int
		loop:
			for foundPeers = 0; foundPeers < tc.expectedPeersFound; foundPeers++ {
				select {
				case <-ctx.Done():
					break loop
				case <-cpeers:
				}
			}

			assert.Equal(t, len(cpeers), 0)
			assert.Equal(t, tc.expectedPeersFound, foundPeers)
		})
	}
}

func TestAnnounceForPeriod(t *testing.T) {
}
