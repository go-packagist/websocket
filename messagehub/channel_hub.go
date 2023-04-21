package messagehub

import "github.com/go-fires/websocket/contracts"

type ChannelHub struct {
	message chan *contracts.Message
}

var _ contracts.MessageHub = (*ChannelHub)(nil)

func NewChannelHub() *ChannelHub {
	return &ChannelHub{
		message: make(chan *contracts.Message, 1024),
	}
}
func (m *ChannelHub) Push(message *contracts.Message) error {
	m.message <- message

	return nil
}

func (m *ChannelHub) Pop() <-chan *contracts.Message {
	return m.message
}
