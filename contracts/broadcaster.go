package contracts

type Broadcaster interface {
	Subscribe(channel Channel, subscriber ...Subscriber)
	Unsubscribe(channel Channel, subscriber ...Subscriber)

	Send(message *Message) // 广播消息
}
