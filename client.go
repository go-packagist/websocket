package websocket

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"sync"
	"time"
)

var (
	ErrClientAlreadySubscribed = errors.New("the client is already subscribed to the channel")
	ErrClientNotSubscribed     = errors.New("the client is not subscribed to the channel")
)

type Client struct {
	id       string
	conn     *websocket.Conn
	channels map[Channel]struct{}

	readChan  chan *Message
	writeChan chan *Message
	errorChan chan error
	closeChan chan struct{}

	mu   sync.Mutex
	once sync.Once
}

var ws = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		id:        uuid.NewV4().String(),
		conn:      conn,
		channels:  make(map[Channel]struct{}),
		readChan:  make(chan *Message, 1024),
		writeChan: make(chan *Message, 1024),
		errorChan: make(chan error, 1),
		closeChan: make(chan struct{}),
	}
}

func UpgradeClient(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Client, error) {
	conn, err := ws.Upgrade(w, r, responseHeader)

	if err != nil {
		return nil, err
	}

	return NewClient(conn), nil
}

func (c *Client) Read() <-chan *Message {
	go c.read()

	return c.readChan
}

func (c *Client) read() {
	for {
		select {
		case <-c.closeChan:
			return
		default:
			messageType, message, err := c.conn.ReadMessage()
			if err != nil {
				c.errorChan <- err
				c.Close()
				return
			}

			var payload *Payload
			if err := json.Unmarshal(message, &payload); err != nil {
				c.errorChan <- err
				return
			}

			c.readChan <- &Message{
				Type:    messageType,
				Payload: payload,
			}
		}
	}
}

func (c *Client) Push(message *Message) {
	c.writeChan <- message
}

func (c *Client) Write() chan<- *Message {
	go c.write()

	return c.writeChan
}

func (c *Client) write() {
	for {
		select {
		case message := <-c.writeChan:
			if payload, err := message.Payload.Marshal(); err != nil {
				c.errorChan <- err
				continue
			} else if err := c.conn.WriteMessage(message.Type, payload); err != nil {
				c.errorChan <- err
				continue
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Client) Closed() <-chan struct{} {
	return c.closeChan
}

func (c *Client) Close() {
	c.once.Do(func() {
		close(c.closeChan)
		_ = c.conn.Close()
	})
}

func (c *Client) Error() <-chan error {
	return c.errorChan
}

func (c *Client) Ping() {
	if err := c.conn.WriteControl(PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
		c.errorChan <- err
		c.Close()
	}
}

func (c *Client) Subscribe(channel Channel) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.In(channel) {
		return ErrClientAlreadySubscribed
	}

	c.channels[channel] = struct{}{}

	return nil
}

func (c *Client) Unsubscribe(channel Channel) error {
	if !c.In(channel) {
		return ErrClientNotSubscribed
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.channels, channel)

	return nil
}

func (c *Client) UnsubscribeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.channels = make(map[Channel]struct{})
}

func (c *Client) In(channel Channel) bool {
	_, ok := c.channels[channel]

	return ok
}

func (c *Client) Channels() []Channel {
	channels := make([]Channel, 0, len(c.channels))

	for channel := range c.channels {
		channels = append(channels, channel)
	}

	return channels
}

func (c *Client) ID() string {
	return c.id
}
