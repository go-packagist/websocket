package websocket

import (
	"encoding/json"
	"github.com/go-fires/websocket/contracts"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
)

type Client struct {
	socketid string
	conn     *websocket.Conn
	closed   chan struct{}

	sendChan    chan *contracts.Message
	receiveChan chan *contracts.Message
}

var _ contracts.Client = (*Client)(nil)

func NewClient(conn *websocket.Conn) (*Client, error) {
	c := &Client{
		conn:        conn,
		sendChan:    make(chan *contracts.Message, 1024),
		receiveChan: make(chan *contracts.Message, 1024),
		closed:      make(chan struct{}),
	}

	c.generateSocketID()

	go c.send()

	return c, nil
}

var ws = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewClientWithHttp(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Client, error) {
	conn, err := ws.Upgrade(w, r, responseHeader)

	if err != nil {
		return nil, err
	}

	return NewClient(conn)
}

func (c *Client) generateSocketID() string {
	return uuid.NewV4().String()
}

func (c *Client) SocketID() string {
	return c.socketid
}

// send to client
func (c *Client) send() {
	for {
		select {
		case payload := <-c.sendChan:
			if err := c.conn.WriteJSON(payload); err != nil {
				log.Printf("Error sending message to client: %v", err)
			}
		case <-c.closed:
			return
		}
	}
}

// Send to client use channel
func (c *Client) Send(message *contracts.Message) {
	c.sendChan <- message
}

func (c *Client) receive() {
	defer func() {
		if err := c.close(); err != nil {
			log.Printf("Error closing client: %v", err)
		}
	}()

	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client: %v", err)
			return
		}

		switch messageType {
		case websocket.CloseMessage:
			return
		case websocket.PingMessage:
			if err := c.conn.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
				log.Printf("Error sending pong message to client: %v", err)
				return
			}
		case websocket.PongMessage:
			continue
		case websocket.TextMessage:
		case websocket.BinaryMessage:
			var payload *contracts.Payload
			if err := json.Unmarshal(message, &payload); err != nil {
				log.Printf("Error unmarshal message from client: %v", err)
				continue // if error, continue to next message
			}

			c.receiveChan <- &contracts.Message{
				MessageType: messageType,
				Payload:     payload,
				From:        c.socketid,
			}
		default:
			log.Printf("Unknown message type: %v", messageType)
		}
	}
}

// Receive from client to channel and return channel
func (c *Client) Receive() <-chan *contracts.Message {
	go c.receive()

	return c.receiveChan
}

// close closes the client.
func (c *Client) close() error {
	close(c.closed)
	return c.conn.Close()
}

// Closed returns a channel that is closed when the client is closed.
func (c *Client) Closed() <-chan struct{} {
	return c.closed
}
