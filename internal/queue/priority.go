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

	messages   []T
	muMessages sync.RWMutex
}

func NewPriorityQueue[T ICounter](name string, tracer MetricsTracer[T]) *PriorityQueue[T] {
	queue := &PriorityQueue[T]{
		name:    name,
		metrics: tracer,

		messages: []T{},
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
	pq.metrics.ItemQueued(pq.name, x.(T))
}

func (pq *PriorityQueue[T]) Pop() (item interface{}) {
	var null T
	if n := len(pq.messages); n > 0 {
		item = pq.messages[n-1]
		pq.metrics.ItemPop(pq.name, item.(T))
		pq.messages, pq.messages[n-1] = pq.messages[:n-1], null
	}
	return item
}
