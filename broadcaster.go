package websocket

import "github.com/go-fires/websocket/contracts"

type Broadcaster struct {
	subscribers map[string][]contracts.Subscriber
}

var _ contracts.Broadcaster = (*Broadcaster)(nil)

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		subscribers: make(map[string][]contracts.Subscriber),
	}
}

func (b *Broadcaster) Broadcast(broadcast contracts.Broadcast) {
	for _, channeler := range broadcast.Channels() {
		for channel, subscribers := range b.subscribers {
			if channel == channeler.Name() {
				for _, subscriber := range subscribers {
					subscriber.Send(&contracts.Payload{
						Channel: channeler.Name(),
						Event:   broadcast.Event(),
						Data:    broadcast.Payload(),
					})
				}
			}
		}
	}
}

func (b *Broadcaster) Subscribe(channel contracts.Channeler, subscriber ...contracts.Subscriber) {
	if _, ok := b.subscribers[channel.Name()]; !ok {
		b.subscribers[channel.Name()] = make([]contracts.Subscriber, 0)
	}

	b.subscribers[channel.Name()] = append(b.subscribers[channel.Name()], subscriber...)
}

func (b *Broadcaster) Unsubscribe(channel contracts.Channeler, subscriber ...contracts.Subscriber) {
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
