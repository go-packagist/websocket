package websocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	assert.Equal(t, "all", AllChannel.Name())
	assert.False(t, AllChannel.IsPrivate())

	assert.Equal(t, "private-foo", NewPrivateChannel("foo").Name())
	assert.True(t, NewPrivateChannel("foo").IsPrivate())

	assert.Equal(t, "presence-foo", NewPresenceChannel("foo").Name())
	assert.True(t, NewPresenceChannel("foo").IsPresence())

	assert.Equal(t, "client-foo", NewClientChannel("foo").Name())
	assert.True(t, NewClientChannel("foo").IsClient())

	assert.False(t, NewChannel("foo").IsPrivate())
	assert.False(t, NewChannel("foo").IsPresence())
	assert.False(t, NewChannel("foo").IsClient())

	assert.True(t, NewChannel("private-foo").IsPrivate())
	assert.False(t, NewChannel("private-foo").IsPresence())
	assert.False(t, NewChannel("private-foo").IsClient())

	assert.False(t, NewChannel("presence-foo").IsPrivate())
	assert.True(t, NewChannel("presence-foo").IsPresence())
	assert.False(t, NewChannel("presence-foo").IsClient())

	assert.False(t, NewChannel("client-foo").IsPrivate())
	assert.False(t, NewChannel("client-foo").IsPresence())
	assert.True(t, NewChannel("client-foo").IsClient())
}
