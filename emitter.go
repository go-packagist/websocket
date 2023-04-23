package websocket

import (
	"sync"
)

type Emitter struct {
	subscribers map[string]Subscriber
	mu          sync.Mutex
}

// NewEmitter create a new emitter
func NewEmitter() *Emitter {
	return &Emitter{
		subscribers: make(map[string]Subscriber, 1024),
	}
}

// Subscribe subscribe a subscriber
func (e *Emitter) Subscribe(subscriber Subscriber) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.subscribers[subscriber.ID()]; !ok {
		e.subscribers[subscriber.ID()] = subscriber
	}
}

// Unsubscribe unsubscribe a subscriber
func (e *Emitter) Unsubscribe(subscriber Subscriber) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.subscribers[subscriber.ID()]; ok {
		delete(e.subscribers, subscriber.ID())
	}
}

// Emit emit a message to all subscribers
func (e *Emitter) Emit(message *Message) {
	if message.Payload == nil {
		return
	}

	if message.Payload.Channel == "" {
		return
	}

	for _, subscriber := range e.subscribers {
		if subscriber.In(message.Payload.Channel) { // check if the subscriber is in the channel
			subscriber.Recv() <- message
		}
	}
}

// Count return the number of subscribers
func (e *Emitter) Count() int {
	return len(e.subscribers)
}
