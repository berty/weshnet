//go:build js

package main

import (
	"fmt"
	"reflect"
	"syscall/js"

	"github.com/libp2p/go-libp2p/core/event"
)

type eventBusFromJS struct {
	helia js.Value
}

// Subscribe creates a new Subscription.
// eventType can be either a pointer to a single event type, or a slice of pointers to
// subscribe to multiple event types at once, under a single subscription (and channel).
//
// Failing to drain the channel may cause publishers to block.
//
// If you want to subscribe to ALL events emitted in the bus, use
// `WildcardSubscription` as the `eventType`:
//
//	eventbus.Subscribe(WildcardSubscription)
//
// Simple example
//
//	sub, err := eventbus.Subscribe(new(EventType))
//	defer sub.Close()
//	for e := range sub.Out() {
//	  event := e.(EventType) // guaranteed safe
//	  [...]
//	}
//
// Multi-type example
//
//	sub, err := eventbus.Subscribe([]interface{}{new(EventA), new(EventB)})
//	defer sub.Close()
//	for e := range sub.Out() {
//	  select e.(type):
//	    case EventA:
//	      [...]
//	    case EventB:
//	      [...]
//	  }
//	}
func (jeb *eventBusFromJS) Subscribe(eventType interface{}, opts ...event.SubscriptionOpt) (event.Subscription, error) {
	fmt.Println("FIXME: ignoring sub to", eventType, opts)
	return nil, nil
}

// Emitter creates a new event emitter.
//
// eventType accepts typed nil pointers, and uses the type information for wiring purposes.
//
// Example:
//
//	em, err := eventbus.Emitter(new(EventT))
//	defer em.Close() // MUST call this after being done with the emitter
//	em.Emit(EventT{})
func (jeb *eventBusFromJS) Emitter(eventType interface{}, opts ...event.EmitterOpt) (event.Emitter, error) {
	panic("not implemented") // TODO: Implement
}

// GetAllEventTypes returns all the event types that this bus knows about
// (having emitters and subscribers). It omits the WildcardSubscription.
//
// The caller is guaranteed that this function will only return value types;
// no pointer types will be returned.
func (jeb *eventBusFromJS) GetAllEventTypes() []reflect.Type {
	panic("not implemented") // TODO: Implement
}
