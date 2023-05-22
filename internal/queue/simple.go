package queue

import (
	"container/list"
	"context"
	"sync"
)

func NewSimpleQueue[T any](name string, tracer MetricsTracer[T]) *SimpleQueue[T] {

	return &SimpleQueue[T]{
		name:    name,
		metrics: tracer,

		items:  list.New(),
		notify: newItemNotify(),
	}
}

// A priorityMessageQueue implements heap.Interface and holds Items.
type SimpleQueue[T any] struct {
	name       string
	items      *list.List
	muMessages sync.RWMutex
	metrics    MetricsTracer[T]

	notify *itemNotify
}

func (q *SimpleQueue[T]) Add(m T) {
	q.muMessages.Lock()
	_ = q.items.PushBack(m)
	q.notify.Broadcast()
	q.metrics.ItemQueued(q.name, m)
	q.muMessages.Unlock()
}

func (q *SimpleQueue[T]) Pop() (m T, ok bool) {
	q.muMessages.Lock()
	if front := q.items.Front(); front != nil {
		m = q.items.Remove(front).(T)
		ok = true
		q.metrics.ItemPop(q.name, m)
	}
	q.muMessages.Unlock()

	return
}

func (q *SimpleQueue[T]) WaitForNewItem(ctx context.Context) (item T, ok bool) {
	for {
		if item, ok = q.Pop(); ok {
			return
		}

		if ok = q.notify.Wait(ctx); !ok {
			return
		}
	}
}

type itemNotify struct {
	signal   chan struct{}
	muSignal sync.Mutex
}

func newItemNotify() *itemNotify {
	return &itemNotify{
		signal: make(chan struct{}, 1),
	}
}

func (m *itemNotify) Wait(ctx context.Context) (ok bool) {
	select {
	case <-m.signal:
		return true
	case <-ctx.Done():
		return false
	}
}

func (m *itemNotify) Broadcast() {
	select {
	case m.signal <- struct{}{}:
	default:
	}
}
