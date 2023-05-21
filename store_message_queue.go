package weshnet

import (
	"github.com/ipfs/go-cid"

	"berty.tech/go-orbit-db/stores/operation"
	"berty.tech/weshnet/internal/queue"
	"berty.tech/weshnet/pkg/protocoltypes"
)

// An Item is something we manage in a priority queue.
type messageItem struct {
	op      operation.Operation
	env     *protocoltypes.MessageEnvelope
	headers *protocoltypes.MessageHeaders
	hash    cid.Cid
}

func (m *messageItem) Counter() uint64 {
	return m.headers.Counter
}

type simpleMessageQueue = queue.SimpleQueue[*messageItem]

func newMessageQueue(name string, tracer queue.MetricsTracer[*messageItem]) *simpleMessageQueue {
	return queue.NewSimpleQueue[*messageItem](name, tracer)
}

type priorityMessageQueue = queue.PriorityQueue[*messageItem]

func newPriorityMessageQueue(name string, tracer queue.MetricsTracer[*messageItem]) *priorityMessageQueue {
	return queue.NewPriorityQueue[*messageItem](name, tracer)
}
