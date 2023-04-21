package websocket

import "github.com/go-fires/websocket/contracts"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) OnOpen(client *Client) {
	client.Send(&contracts.Message{
		Channel: "socketing:connection",
		Event:   "connected",
		Payload: struct {
			Socketid string
		}{
			Socketid: client.SocketID(),
		},
	})
}

func (h *Handler) OnClose(client *Client) {
	defer client.close()

	// unbind
}

func (h *Handler) OnMessage(client *Client, message *contracts.Message) {
}

func (h *Handler) OnError(client *Client, err error) {

}
