package queue

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testSimpleQueue = *SimpleQueue[int]

func newTestSimpleQueue() testSimpleQueue {
	return NewSimpleQueue[int]("test", &noopTracer[int]{})
}

func TestQueue(t *testing.T) {
	queue := newTestSimpleQueue()

	e, ok := queue.Pop()
	require.Equal(t, 0, e)
	require.False(t, ok)

	queue.Add(1)
	e, ok = queue.Pop()
	require.Equal(t, 1, e)
	require.True(t, ok)

	e, ok = queue.Pop()
	require.Equal(t, 0, e)
	require.False(t, ok)
}

func TestSyncQueue(t *testing.T) {
	cases := []struct{ N int }{
		{1}, {10}, {100}, {1000}, {10000},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%d_elements", tc.N)
		t.Run(name, func(t *testing.T) {
			queue := newTestSimpleQueue()

			for i := 0; i < tc.N; i++ {
				queue.Add(i + 1)
			}

			for i := 0; i < tc.N; i++ {
				e, ok := queue.Pop()
				require.Equal(t, i+1, e)
				require.True(t, ok)
			}
		})
	}
}

func TestAsyncQueue(t *testing.T) {
	cases := []struct{ N int }{
		{1}, {10}, {100}, {1000}, {10000},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%d_elements", tc.N)
		t.Run(name, func(t *testing.T) {
			queue := newTestSimpleQueue()

			wg := sync.WaitGroup{}

			wg.Add(tc.N)
			elems := map[int]struct{}{}
			for i := 0; i < tc.N; i++ {
				elems[i+1] = struct{}{}
				go func(i int) {
					queue.Add(i + 1)
					wg.Done()
				}(i)
			}

			wg.Wait()

			for i := 0; i < tc.N; i++ {
				e, ok := queue.Pop()
				require.True(t, ok)

				_, exist := elems[e]
				require.True(t, exist)
				delete(elems, e)
			}

			require.Len(t, elems, 0)
		})
	}
}

func TestWaitnForItemQueue(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cases := []struct{ N int }{
		{1}, {10}, {100}, {1000}, {10000},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%d_elements", tc.N)
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			queue := newTestSimpleQueue()

			cc := make(chan int, tc.N)
			go func() {
				defer close(cc)
				for {
					e, ok := queue.WaitForItem(ctx)
					if !ok {
						return
					}

					cc <- e
				}
			}()

			for i := 0; i < tc.N; i++ {
				queue.Add(i + 1)
			}

			for i := 0; i < tc.N; i++ {
				select {
				case e := <-cc:
					require.Equal(t, i+1, e)
				case <-time.After(time.Second):
					require.FailNow(t, "timeout while waiting for event")
				}
			}
		})
	}
}
