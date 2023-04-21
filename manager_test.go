package websocket

import (
	"testing"

	"github.com/go-fires/websocket/messagehub"
	"github.com/gorilla/websocket"
)

func TestManager(t *testing.T) {
	m := NewManager(messagehub.NewChannelHub())

	client, err := NewClient(&websocket.Conn{})
	if err != nil {
		t.Fatal(err)
	}

	m.Handle(client)
}
