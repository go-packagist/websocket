package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type SubscriberManager struct {
	clients map[string]*Client

	mu sync.Mutex
}

func NewSubscriberManager() *SubscriberManager {
	return &SubscriberManager{
		clients: make(map[string]*Client),
	}
}

func (m *SubscriberManager) NewClient(conn *websocket.Conn) *Client {
	client := NewClient(conn)

	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[client.socketid] = client

	return client
}

func (m *SubscriberManager) NewClientWithHttp(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Client, error) {
	client, err := NewClientWithHttp(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	return m.NewClient(client.conn)
}

// Handle client
func (m *SubscriberManager) Handle(client *Client) {
	for {
		select {
		case msg := <-client.Read():
			// fmt.Println("Handle read from client:", string(msg.Data))
			m.kernelMesssage <- &Message{
				From: client.SocketID(),
				Type: msg.MessageType,
				Data: msg.Data,
				To:   []string{"all"},
			}
		case <-client.Closed():
			m.RemoveClient(client.SocketID())
			return
		}
	}
}
