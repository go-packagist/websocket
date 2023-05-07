package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

var (
	TextMessage   = websocket.TextMessage
	BinaryMessage = websocket.BinaryMessage
	PingMessage   = websocket.PingMessage
	PongMessage   = websocket.PongMessage
	CloseMessage  = websocket.CloseMessage
)

type Message struct {
	Type    int      `json:"type"`
	Payload *Payload `json:"payload"`
}

type Payload struct {
	Channel Channel     `json:"channel,omitempty"`
	Event   string      `json:"event,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (p *Payload) Unmarshal(payload []byte) error {
	return json.Unmarshal(payload, &p)
}

func (p *Payload) Marshal() ([]byte, error) {
	return json.Marshal(p)
}
