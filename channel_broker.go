package websocket

type ChannelBroker struct {
	message chan *Message
}

func NewChannelBroker() *ChannelBroker {
	return &ChannelBroker{
		message: make(chan *Message, 1024),
	}
}

func (c *ChannelBroker) Publish() chan<- *Message {
	return c.message
}

func (c *ChannelBroker) Consume() <-chan *Message {
	return c.message
}

func (c *ChannelBroker) Close() error {
	close(c.message)

	return nil
}
