package websocket

type Handler interface {
	Opened(socketer Socketable)
	Closed(socketer Socketable)
	Message(socketer Socketable, message *Message)
	Error(socketer Socketable, err error)
}
