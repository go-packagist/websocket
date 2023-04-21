package contracts

type Sender interface {
	Send(message *Message)
}

type Receiver interface {
	Receive() <-chan *Message
}

type Publisher interface {
	Receiver
}

type Subscriber interface {
	Sender

	Subscribe() <-chan Channel
	Unsubscribe() <-chan Channel
}

type Client interface {
	Subscriber
	Publisher
}
