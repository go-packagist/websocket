package websocket

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient(t *testing.T) {
	c := NewClient(nil)

	assert.False(t, c.In("not-exist"))
	assert.Equal(t, 0, len(c.Channels()))

	assert.NoError(t, c.Subscribe("channel-1"))
	assert.Equal(t, ErrClientAlreadySubscribed, c.Subscribe("channel-1"))
}
