package main

import (
	"fmt"
	"net/http"

	"github.com/go-fires/websocket"
)

type Handler struct {
}

var _ websocket.Handler = (*Handler)(nil)

// Opened is called when a new socket is opened
func (w *Handler) Opened(socketer websocket.Socketable) {
	fmt.Printf("%s is connected\n", socketer.ID())
}

// Closed is called when a socket is closed
func (w *Handler) Closed(socketer websocket.Socketable) {
	fmt.Printf("%s is disconnected\n", socketer.ID())
}

// Message is called when a socket sends a message
func (w *Handler) Message(socketer websocket.Socketable, message *websocket.Message) {
	if message == nil {
		return
	}

	socketer.Push(message)
	fmt.Println("message", message)
}

func (w *Handler) Error(socketer websocket.Socketable, err error) {
	fmt.Printf("Error[%s]: %s\n", socketer.ID(), err.Error())
}

var engine = websocket.NewEngine(websocket.NewChannelHub(), new(Handler))

func main() {
	http.HandleFunc("/wss", func(w http.ResponseWriter, r *http.Request) {
		socket, _ := websocket.NewSocketWithHttp(w, r, nil)

		engine.Handle(socket)
	})

	err := http.ListenAndServe(":1027", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
