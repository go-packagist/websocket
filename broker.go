package websocket

type Broker interface {
	Publish() chan<- *Message
	Consume() <-chan *Message
}
