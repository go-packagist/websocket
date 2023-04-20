package contracts

type Broadcaster interface {
	Subscribe(channel Channeler, subscriber ...Subscriber)
	Unsubscribe(channel Channeler, subscriber ...Subscriber)

	Broadcast(broadcast Broadcast)
}
