package queue

import (
	"container/heap"
	"sync"
)

type ICounter interface {
	Counter() uint64
}

// A priorityMessageQueue implements heap.Interface and holds Items.
type PriorityQueue[T ICounter] struct {
	name    string
	metrics MetricsTracer[T]

	items      []T
	muMessages sync.RWMutex
}

func NewPriorityQueue[T ICounter](name string, tracer MetricsTracer[T]) *PriorityQueue[T] {
	queue := &PriorityQueue[T]{
		name:    name,
		metrics: tracer,

		items: []T{},
	}

	heap.Init(queue)
	return queue
}

func (pq *PriorityQueue[T]) Add(m T) {
	pq.muMessages.Lock()
	heap.Push(pq, m)
	pq.metrics.ItemQueued(pq.name, m)
	pq.muMessages.Unlock()
}

func (pq *PriorityQueue[T]) NextAll(cb func(next T) error) error {
	pq.muMessages.Lock()
	defer pq.muMessages.Unlock()

	for len(pq.items) > 0 {
		item := heap.Pop(pq).(T)
		if err := cb(item); err != nil {
			return err
		}
	}

	return nil
}

func (pq *PriorityQueue[T]) Next() (item T) {
	pq.muMessages.Lock()
	if len(pq.items) > 0 {
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
	l = len(pq.items)
	return
}

func (pq *PriorityQueue[T]) Less(i, j int) bool {
	// We want Pop to give us the lowest, not highest, priority so we use lower than here.
	return pq.items[i].Counter() < pq.items[j].Counter()
}

func (pq *PriorityQueue[T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

func (pq *PriorityQueue[T]) Push(x interface{}) {
	pq.items = append(pq.items, x.(T))
	pq.metrics.ItemQueued(pq.name, x.(T))
}

func (pq *PriorityQueue[T]) Pop() (item interface{}) {
	var null T
	if n := len(pq.items); n > 0 {
		item = pq.items[n-1]
		pq.metrics.ItemPop(pq.name, item.(T))
		pq.items, pq.items[n-1] = pq.items[:n-1], null
	}
	return item
}
