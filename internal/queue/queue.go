package queue

import (
	"container/heap"
	"container/list"
	"context"
	"sync"
)

type ICounter interface {
	Counter() uint64
}

// A priorityMessageQueue implements heap.Interface and holds Items.
type PriorityQueue[T ICounter] struct {
	messages   []T
	muMessages sync.RWMutex
}

func NewPriorityQueue[T ICounter]() *PriorityQueue[T] {
	queue := &PriorityQueue[T]{
		messages: []T{},
	}

	heap.Init(queue)
	return queue
}

func (pq *PriorityQueue[T]) Add(m T) {
	pq.muMessages.Lock()
	heap.Push(pq, m)
	pq.muMessages.Unlock()
}

func (pq *PriorityQueue[T]) Next() (item T) {
	pq.muMessages.Lock()
	if len(pq.messages) > 0 {
		item = heap.Pop(pq).(T)
	}
	pq.muMessages.Unlock()
	return
}

func (pq *PriorityQueue[T]) Size() (l int) {
	pq.muMessages.RLock()
	l = pq.Len()
	pq.muMessages.RUnlock()
	return
}

func (pq *PriorityQueue[T]) Len() (l int) {
	l = len(pq.messages)
	return
}

func (pq *PriorityQueue[T]) Less(i, j int) bool {
	// We want Pop to give us the lowest, not highest, priority so we use lower than here.
	return pq.messages[i].Counter() < pq.messages[j].Counter()
}

func (pq *PriorityQueue[T]) Swap(i, j int) {
	pq.messages[i], pq.messages[j] = pq.messages[j], pq.messages[i]
}

func (pq *PriorityQueue[T]) Push(x interface{}) {
	pq.messages = append(pq.messages, x.(T))
}

func (pq *PriorityQueue[T]) Pop() (item interface{}) {
	var null T
	if n := len(pq.messages); n > 0 {
		item = pq.messages[n-1]
		pq.messages, pq.messages[n-1] = pq.messages[:n-1], null
	}
	return item
}

func NewSimpleQueue[T any]() *SimpleQueue[T] {
	return &SimpleQueue[T]{
		items:  list.New(),
		notify: newItemNotify(),
	}
}

// A priorityMessageQueue implements heap.Interface and holds Items.
type SimpleQueue[T any] struct {
	items      *list.List
	muMessages sync.RWMutex

	notify *itemNotify
}

func (q *SimpleQueue[T]) Add(m T) {
	q.muMessages.Lock()
	_ = q.items.PushBack(m)
	q.notify.Broadcast()
	q.muMessages.Unlock()
}

func (q *SimpleQueue[T]) Pop() (m T, ok bool) {
	q.muMessages.Lock()
	if front := q.items.Front(); front != nil {
		m = q.items.Remove(front).(T)
		ok = true
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
