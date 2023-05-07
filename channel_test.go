package websocket

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChannel(t *testing.T) {
	assert.Equal(t, "channel", Channel("channel").String())
	assert.Equal(t, "channel", NewChannel("channel").String())
	assert.Equal(t, PrivatePrefix+"channel", NewPrivateChannel("channel").Name())
	assert.Equal(t, PrivateEncryptedPrefix+"channel", NewPrivateEncryptedChannel("channel").Name())
	assert.Equal(t, PresencePrefix+"channel", NewPresenceChannel("channel").Name())
	assert.True(t, NewPrivateChannel("channel").IsPrivate())
	assert.True(t, NewPrivateEncryptedChannel("channel").IsPrivateEncrypted())
	assert.True(t, NewPrivateEncryptedChannel("channel").IsPrivate())
	assert.True(t, NewPresenceChannel("channel").IsPresence())

	assert.False(t, NewChannel("channel").IsPrivate())
	assert.False(t, NewChannel("channel").IsPrivateEncrypted())
	assert.False(t, NewChannel("channel").IsPresence())
}
