package contracts

type Message struct {
	MessageType int         `json:"-"` // websocket message type see https://godoc.org/github.com/gorilla/websocket#pkg-constants
	Owner       string      `json:"-"` // sender client socketid
	Channel     string      `json:"channel"`
	Event       string      `json:"event"`
	Payload     interface{} `json:"payload"`
}

// MessageHub is a message hub interface, for example, you can use redis, mq... as a message hub.
type MessageHub interface {
	In() chan<- *Message
	Out() <-chan *Message
}
