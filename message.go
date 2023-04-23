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
	Channel string      `json:"channel"`
	Event   string      `json:"event"`
	Data    interface{} `json:"data"`
}

func (p *Payload) Unmarshal(payload []byte) error {
	return json.Unmarshal(payload, &p)
}

func (p *Payload) Marshal() ([]byte, error) {
	return json.Marshal(p)
}
