package websocket

type Hub interface {
	Recv() chan<- *Message
	Read() <-chan *Message
}
