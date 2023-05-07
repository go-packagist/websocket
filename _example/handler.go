package main

import (
	"fmt"
	"github.com/go-fires/websocket"
)

type Handler struct {
}

var _ websocket.Handler = (*Handler)(nil)

func (h *Handler) Opened(client *websocket.Client) {
	fmt.Println("Opened", client.ID())
}

func (h *Handler) Message(client *websocket.Client, message *websocket.Message) {
	client.Subscribe(message.Payload.Channel)
	fmt.Println("Message", client.ID(), message)
	server.Broadcast(message)
}

func (h *Handler) Error(client *websocket.Client, err error) {
	fmt.Println("Error", client.ID(), err)
}

func (h *Handler) Closed(client *websocket.Client) {
	fmt.Println("Closed", client.ID())
}
