package websocket

type ChannelHub struct {
	message chan *Message
}

var _ Hub = (*ChannelHub)(nil)

func NewChannelHub() *ChannelHub {
	return &ChannelHub{
		message: make(chan *Message, 1024),
	}
}

func (c *ChannelHub) Recv() chan<- *Message {
	return c.message
}

func (c *ChannelHub) Read() <-chan *Message {
	return c.message
}
