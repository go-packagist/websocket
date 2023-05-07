package websocket

type NullHandler struct {
}

var _ Handler = (*NullHandler)(nil)

func (n *NullHandler) Opened(client *Client) {
}

func (n *NullHandler) Message(client *Client, message *Message) {
}

func (n *NullHandler) Error(client *Client, err error) {
}

func (n *NullHandler) Closed(client *Client) {
}
