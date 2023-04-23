package websocket

type Subscriber interface {
	ID() string
	Recv() chan<- *Message

	Join(channel string) error
	Leave(channel string) error
	LeaveAll() error
	In(channel string) bool
}

type Publisher interface {
	ID() string
	Read() <-chan *Message
	Push(message *Message)
}

type Socketable interface {
	Subscriber
	Publisher

	Closed() <-chan struct{}
	Error() <-chan error
}
