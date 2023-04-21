package v1

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// Message websocket message
type Message struct {
	MessageType int
	Message     []byte
}

var ws = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Socket struct {
	id   string
	conn *websocket.Conn
}

func NewSocket(conn *websocket.Conn) *Socket {
	return &Socket{
		id:   uuid.NewV4().String(),
		conn: conn,
	}
}

// NewSocketWithHttp create socket with http
func NewSocketWithHttp(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Socket, error) {
	conn, err := ws.Upgrade(w, r, responseHeader)

	if err != nil {
		return nil, err
	}

	return NewSocket(conn), nil
}

// Close  socket
func (s *Socket) Close() {
	_ = s.conn.Close()
}

// ID return socket id
func (s *Socket) ID() string {
	return s.id
}

// State return send, recv, quit channel
func (s *Socket) State() (chan<- *Message, <-chan *Message, <-chan struct{}) {
	send := make(chan *Message, 1024)
	recv := make(chan *Message, 1024)
	closed := make(chan struct{}, 1)
	quit := make(chan struct{}, 0)

	// close to quit
	go func() {
		for {
			select {
			case <-closed:
				quit <- struct{}{}
				return
			}
		}
	}()

	// recv request
	go func() {
		for {
			select {
			case <-closed:
				return
			default:
				messageType, msg, err := s.conn.ReadMessage()
				if err != nil {
					closed <- struct{}{}
					return
				}

				recv <- &Message{
					MessageType: messageType,
					Message:     msg,
				}
			}
		}
	}()

	// send response
	go func() {
		for {
			select {
			case <-closed:
				return
			default:
				msg := <-send

				_ = s.conn.WriteMessage(msg.MessageType, msg.Message)
			}
		}
	}()

	// send ping to control connection
	go func() {
		for {
			select {
			case <-time.Tick(25 * time.Second):
				err := s.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second))
				if err != nil {
					closed <- struct{}{}
					return
				}
			}
		}
	}()

	return send, recv, quit
}
