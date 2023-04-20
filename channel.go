package websocket

import "github.com/go-fires/websocket/contracts"

type Channel struct {
	name string
}

var _ contracts.Channeler = (*Channel)(nil)

func NewChannel(name string) *Channel {
	return &Channel{
		name: name,
	}
}

func NewPrivateChannel(name string) *Channel {
	return NewChannel("private-" + name)
}

func NewPresenceChannel(name string) *Channel {
	return NewChannel("presence-" + name)
}

func NewClientChannel(name string) *Channel {
	return NewChannel("client-" + name)
}

func (c *Channel) Name() string {
	return c.name
}

func (c *Channel) IsPrivate() bool {
	return c.name[:8] == "private-"
}

func (c *Channel) IsPresence() bool {
	return c.name[:10] == "presence-"
}
