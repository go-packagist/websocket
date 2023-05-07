package main

import (
	"fmt"
	"github.com/go-fires/websocket"
	"net/http"
)

var server = websocket.NewServer(
	&Handler{}, websocket.NewChannelBroker())

func main() {
	go server.Start()

	http.HandleFunc("/wss", func(w http.ResponseWriter, r *http.Request) {
		client, err := websocket.UpgradeClient(w, r, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		server.Handle(client)
	})

	err := http.ListenAndServe("127.0.0.1:1026", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
