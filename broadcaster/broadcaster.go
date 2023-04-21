package broadcaster

import (
	"sync"

	"github.com/go-fires/websocket/contracts"
)

type Broadcaster struct {
	subscribers map[string][]contracts.Subscriber // channel name -> subscribers

	mu sync.Mutex
}

var _ contracts.Broadcaster = (*Broadcaster)(nil)

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		subscribers: make(map[string][]contracts.Subscriber),
	}
}

func (b *Broadcaster) Send(message *contracts.Message) {
	if subscribers, ok := b.subscribers[message.Channel]; !ok {
		return
	} else {
		for _, subscriber := range subscribers {
			subscriber.Receive() <- message
		}
	}
}

func (b *Broadcaster) Subscribe(channel contracts.Channel, subscriber ...contracts.Subscriber) {
	if _, ok := b.subscribers[channel.Name()]; !ok {
		b.mu.Lock()
		b.subscribers[channel.Name()] = make([]contracts.Subscriber, 0)
		b.mu.Unlock()
	}

	b.subscribers[channel.Name()] = append(b.subscribers[channel.Name()], subscriber...)
}

func (b *Broadcaster) Unsubscribe(channel contracts.Channel, subscriber ...contracts.Subscriber) {
	if _, ok := b.subscribers[channel.Name()]; !ok {
		return
	}

	for _, s := range subscriber {
		for i, sub := range b.subscribers[channel.Name()] {
			if sub == s {
				b.subscribers[channel.Name()] = append(b.subscribers[channel.Name()][:i], b.subscribers[channel.Name()][i+1:]...)
			}
		}
	}
}
