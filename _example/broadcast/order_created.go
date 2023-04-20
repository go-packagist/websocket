package broadcast

import (
	"github.com/go-fires/websocket"
	"github.com/go-fires/websocket/contracts"
)

type OrderCreated struct {
	OrderId     int    `json:"order_id"`
	OrderStatus string `json:"order_status"`
	CreatedAt   string `json:"created_at"`
}

var _ contracts.Broadcast = (*OrderCreated)(nil)

func (b *OrderCreated) Channels() []contracts.Channeler {
	return []contracts.Channeler{
		websocket.NewChannel("order-created"),
	}
}

func (b *OrderCreated) Event() string {
	return "order-created"
}

func (b *OrderCreated) Payload() interface{} {
	return b
}
