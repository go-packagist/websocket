package broadcast

import (
	"github.com/go-fires/websocket"
	"github.com/go-fires/websocket/contracts"
)

type OrderCreating struct {
	OrderId     int    `json:"order_id"`
	OrderStatus string `json:"order_status"`
	CreatedAt   string `json:"created_at"`
}

var _ contracts.Broadcast = (*OrderCreating)(nil)

func (b *OrderCreating) Channels() []contracts.Channeler {
	return []contracts.Channeler{
		websocket.NewChannel("order-creating"),
	}
}

func (b *OrderCreating) Event() string {
	return "order-created"
}

func (b *OrderCreating) Payload() interface{} {
	return b
}
