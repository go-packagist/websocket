package contracts

type WebsocketMessage struct {
}

type WebsocketClient interface {
	Subscriber

	Receive() (*Payload, error)
	Close() error
}
