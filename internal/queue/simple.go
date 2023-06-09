package queue

import (
	"container/list"
	"context"
	"sync"
)

type SimpleQueue[T any] struct {
	name    string
	list    *list.List
	metrics MetricsTracer[T]
	signal  chan struct{}
	mu      sync.Mutex
}

func NewSimpleQueue[T any](name string, tracer MetricsTracer[T]) *SimpleQueue[T] {
	return &SimpleQueue[T]{
		name:    name,
		metrics: tracer,
		list:    list.New(),
		signal:  make(chan struct{}),
	}
}

// Add pushes an item to the queue
func (q *SimpleQueue[T]) Add(m T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	_ = q.list.PushBack(m)
	q.metrics.ItemQueued(q.name, m)

	// signal that we got a new item
	select {
	case q.signal <- struct{}{}:
	default:
	}
}

// Pop removes and returns the first item from the queue.
// If the queue is empty, the second returned value will be false.
func (q *SimpleQueue[T]) Pop() (m T, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.list.Len() > 0 {
		element := q.list.Front()
		q.list.Remove(element)
		m = element.Value.(T)
		ok = true
	}

	return
}

// WaitForItem blocks until a new item is available or the context is canceled.
// It returns the new item along with a boolean value indicating whether context has expired.
func (q *SimpleQueue[T]) WaitForItem(ctx context.Context) (item T, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Keep attempting to retrieve a new item until the context is canceled
	for ctx.Err() == nil {
		if q.list.Len() == 0 {
			// queue is empty, wait for either a signal of a new item or a
			// context cancellation
			q.mu.Unlock()
			select {
			case <-q.signal:
			case <-ctx.Done():
			}
			q.mu.Lock()

			continue
		}

		// pop front item from the queue
		element := q.list.Front()
		q.list.Remove(element)

		return element.Value.(T), true
	}

	return
}
