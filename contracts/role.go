package contracts

type Publisher interface {
	Client
}

type Subscriber interface {
	Client
}

type Sender interface {
	Send(message *Message)
}

type Receiver interface {
	Receive() <-chan *Message
}

type Client interface {
	Sender
	Receiver
}
