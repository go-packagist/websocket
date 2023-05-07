package websocket

type Handler interface {
	Opened(client *Client)
	Message(client *Client, message *Message)
	Error(client *Client, err error)
	Closed(client *Client)
}
