package websocket

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var ws = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Socketer struct {
	id       string
	conn     *websocket.Conn
	recvChan chan *Message
	readChan chan *Message

	opened chan struct{}
	closed chan struct{}
	errors chan error
	once   sync.Once

	channels map[string]struct{} // subscribed channels
	mu       sync.Mutex
}

var _ Socketable = (*Socketer)(nil)

func NewSocket(conn *websocket.Conn) *Socketer {
	s := &Socketer{
		id:       uuid.NewV4().String(),
		conn:     conn,
		recvChan: make(chan *Message, 1024),
		readChan: make(chan *Message, 1024),
		errors:   make(chan error, 1),
		closed:   make(chan struct{}),
		channels: make(map[string]struct{}, 0),
	}

	go s.recv() // recv message from emitter
	go s.read() // read message from client
	go s.ping() // send ping to control connection

	return s
}

// NewSocketWithHttp create socket with http
func NewSocketWithHttp(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Socketer, error) {
	conn, err := ws.Upgrade(w, r, responseHeader)

	if err != nil {
		return nil, err
	}

	return NewSocket(conn), nil
}

// read message from client and send to recvChan
func (s *Socketer) read() {
	for {
		select {
		case <-s.closed:
			return
		default:
			messageType, data, err := s.conn.ReadMessage()
			if err != nil {
				s.Close()
				return
			}

			var payload Payload
			if err := json.Unmarshal(data, &payload); err != nil {
				s.errors <- err
				// log.Printf("unmarshal payload error: %v, data: %s", err, string(data))
				continue
			}

			s.readChan <- &Message{
				Type:    messageType,
				Payload: &payload,
			}
		}
	}
}

// recv message from Emitter and send to client
func (s *Socketer) recv() {
	for {
		select {
		case <-s.closed:
			return
		case msg := <-s.recvChan:
			if msg == nil {
				continue
			}

			if payload, err := msg.Payload.Marshal(); err != nil {
				s.errors <- err
				// log.Printf("marshal payload error: %v, payload: %+v", err, msg.Payload)
			} else {
				if err := s.conn.WriteMessage(msg.Type, payload); err != nil {
					s.errors <- err
					// log.Printf("write message error: %v, message: %+v", err, msg)
				}
			}
		}
	}
}

// ping send ping to control connection
func (s *Socketer) ping() {
	for {
		select {
		case <-s.closed:
			return
		case <-time.Tick(25 * time.Second):
			err := s.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second))
			if err != nil {
				s.errors <- err
				s.Close()
				return
			}
		}
	}
}

// Closed return closed channel
func (s *Socketer) Closed() <-chan struct{} {
	return s.closed
}

func (s *Socketer) Error() <-chan error {
	return s.errors
}

// Close socket
func (s *Socketer) Close() {
	s.once.Do(func() {
		close(s.recvChan)
		close(s.readChan)
		close(s.closed)
		_ = s.conn.Close()
	})
}

// ID return socket id
func (s *Socketer) ID() string {
	return s.id
}

// Read return read from client channel
func (s *Socketer) Read() <-chan *Message {
	return s.readChan
}

// Recv return recv from emitter channel
func (s *Socketer) Recv() chan<- *Message {
	return s.recvChan
}

// Push message to recv channel
func (s *Socketer) Push(message *Message) {
	s.recvChan <- message
}

func (s *Socketer) Join(channel string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if channel == "" {
		return errors.New("channel is empty")
	}

	if _, ok := s.channels[channel]; ok {
		return errors.New("already joined")
	}

	s.channels[channel] = struct{}{}

	return nil
}

func (s *Socketer) Leave(channel string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.channels[channel]; ok {
		delete(s.channels, channel)
		return nil
	}

	return errors.New("not joined")
}

func (s *Socketer) LeaveAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.channels = make(map[string]struct{})

	return nil
}

func (s *Socketer) In(channel string) bool {
	_, ok := s.channels[channel]

	return ok
}

func (s *Socketer) Channels() []string {
	channels := make([]string, 0, len(s.channels))

	for channel := range s.channels {
		channels = append(channels, channel)
	}

	return channels
}
