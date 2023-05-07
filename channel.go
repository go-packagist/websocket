package websocket

import "strings"

const (
	PrivatePrefix          = "private-"
	PrivateEncryptedPrefix = "private-encrypted-"
	PresencePrefix         = "presence-"
)

// Channel represents a channel to websocket message.
type Channel string

// NewChannel returns a new channel.
func NewChannel(name string) Channel {
	return Channel(name)
}

// NewPrivateChannel returns a new private channel.
func NewPrivateChannel(name string) Channel {
	return NewChannel(PrivatePrefix + name)
}

// NewPrivateEncryptedChannel returns a new private encrypted channel.
func NewPrivateEncryptedChannel(name string) Channel {
	return NewChannel(PrivateEncryptedPrefix + name)
}

// NewPresenceChannel returns a new presence channel.
func NewPresenceChannel(name string) Channel {
	return NewChannel(PresencePrefix + name)
}

// String returns the string representation of the channel.
func (c Channel) String() string {
	return string(c)
}

// Name is alias of String.
func (c Channel) Name() string {
	return c.String()
}

// IsPrivate returns true if the channel is a private channel.
func (c Channel) IsPrivate() bool {
	return strings.Index(c.String(), PrivatePrefix) == 0
}

// IsPrivateEncrypted returns true if the channel is a private encrypted channel.
func (c Channel) IsPrivateEncrypted() bool {
	return strings.Index(c.String(), PrivateEncryptedPrefix) == 0
}

// IsPresence returns true if the channel is a presence channel.
func (c Channel) IsPresence() bool {
	return strings.Index(c.String(), PresencePrefix) == 0
}
