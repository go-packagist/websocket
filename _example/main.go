package main

import (
	"fmt"
	"github.com/go-fires/websocket"
	"github.com/go-fires/websocket/_example/broadcast"
	"github.com/go-fires/websocket/contracts"
)

type subscriber struct {
}

var _ contracts.Subscriber = (*subscriber)(nil)

func (s *subscriber) Send(payload *contracts.Payload) {
	fmt.Println(payload)
}

func main() {
	broadcaster := websocket.NewBroadcaster()

	broadcaster.Subscribe(websocket.NewChannel("order-created"), &subscriber{})
	broadcaster.Subscribe(websocket.NewChannel("order-created"), &subscriber{})
	broadcaster.Subscribe(websocket.NewChannel("order-creating"), &subscriber{})

	broadcaster.Broadcast(&broadcast.OrderCreated{
		OrderId:     1,
		OrderStatus: "created",
		CreatedAt:   "2020-01-01 00:00:00",
	})

	broadcaster.Broadcast(&broadcast.OrderCreating{
		OrderId:     2,
		OrderStatus: "creating",
		CreatedAt:   "2020-01-01 00:00:00",
	})
}
