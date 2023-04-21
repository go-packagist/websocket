package websocket

import (
	"strings"

	"github.com/go-fires/websocket/contracts"
)

const (
	PrivateChannelPrefix  = "private-"
	PresenceChannelPrefix = "presence-"
	ClientChannelPrefix   = "client-"
)

var (
	AllChannel = NewChannel("all")
)

type Channel struct {
	name string
}

var _ contracts.Channel = (*Channel)(nil)

func NewChannel(name string) *Channel {
	return &Channel{
		name: name,
	}
}

func NewPrivateChannel(name string) *Channel {
	return NewChannel(PrivateChannelPrefix + name)
}

func NewPresenceChannel(name string) *Channel {
	return NewChannel(PresenceChannelPrefix + name)
}

func NewClientChannel(name string) *Channel {
	return NewChannel(ClientChannelPrefix + name)
}

func (c *Channel) Name() string {
	return c.name
}

func (c *Channel) IsPrivate() bool {
	return strings.Index(c.name, PrivateChannelPrefix) == 0
}

func (c *Channel) IsPresence() bool {
	return strings.Index(c.name, PresenceChannelPrefix) == 0
}

func (c *Channel) IsClient() bool {
	return strings.Index(c.name, ClientChannelPrefix) == 0
}
