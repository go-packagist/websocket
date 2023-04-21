package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-fires/websocket/contracts"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type Client struct {
	socketid string
	conn     *websocket.Conn

	sendChan        chan *contracts.Message
	receiveChan     chan *contracts.Message
	subscribeChan   chan contracts.Channel
	unsubscribeChan chan contracts.Channel
	closed          chan struct{}
}

var _ contracts.Client = (*Client)(nil)

func NewClient(conn *websocket.Conn) (*Client, error) {
	c := &Client{
		conn:            conn,
		sendChan:        make(chan *contracts.Message, 1024),
		receiveChan:     make(chan *contracts.Message, 1024),
		subscribeChan:   make(chan contracts.Channel, 1),
		unsubscribeChan: make(chan contracts.Channel, 1),
		closed:          make(chan struct{}),
	}

	c.generateSocketID()

	go c.write()

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

// write() to client
func (c *Client) write() {
	for {
		select {
		case payload := <-c.receiveChan:
			if err := c.conn.WriteJSON(payload); err != nil {
				log.Printf("Error sending message to client: %v", err)
			}
		case <-c.closed:
			return
		}
	}
}

// Send client to server
func (c *Client) Send() <-chan *contracts.Message {
	go c.send()

	return c.sendChan
}

func (c *Client) send() {
	defer func() {
		if err := c.close(); err != nil {
			log.Printf("Error closing client: %v", err)
		}
	}()

	for {
		messageType, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client: %v", err)
			return
		}

		switch messageType {
		case websocket.CloseMessage:
			return
		case websocket.PingMessage:
			c.sendChan <- &contracts.Message{
				MessageType: websocket.PingMessage,
				Owner:       c.socketid,
				Event:       "socket:ping",
				Payload:     msg,
			}
		case websocket.PongMessage:
			c.sendChan <- &contracts.Message{
				MessageType: websocket.PongMessage,
				Owner:       c.socketid,
				Event:       "socket:pong",
				Payload:     msg,
			}
		case websocket.TextMessage:
		case websocket.BinaryMessage:
			var message *contracts.Message
			if err := json.Unmarshal(msg, &message); err != nil {
				log.Printf("Error unmarshal message from client: %v", err)
				continue // if error, continue to next message
			}

			c.sendChan <- &contracts.Message{
				MessageType: messageType,
				Owner:       c.socketid,
				Channel:     message.Channel,
				Event:       message.Event,
				Payload:     message.Payload,
			}
		default:
			log.Printf("Unknown message type: %v", messageType)
		}
	}
}

// Receive server to client
func (c *Client) Receive() chan<- *contracts.Message {
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

func (c *Client) Subscribe() <-chan contracts.Channel {
	return c.subscribeChan
}

func (c *Client) Unsubscribe() <-chan contracts.Channel {
	return c.unsubscribeChan
}
