package websocket

import (
	"sync"

	"github.com/go-fires/websocket/contracts"
)

type Manager struct {
	clients     map[string]*Client
	messageHub  contracts.MessageHub
	broadcaster contracts.Broadcaster

	mu sync.Mutex
}

func NewManager(messageHub contracts.MessageHub) *Manager {
	return &Manager{
		clients:    make(map[string]*Client, 1024),
		messageHub: messageHub,
		// broadcaster: broadcaster,
	}
}

func (m *Manager) AddClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[client.SocketID()] = client
}

func (m *Manager) RemoveClient(socketid string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.clients, socketid)
}

func (m *Manager) GetClient(socketid string) *Client {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.clients[socketid]
}

func (m *Manager) Clients() map[string]*Client {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.clients
}

func (m *Manager) ClientCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.clients)
}

func (m *Manager) Broadcast(msg *contracts.Message) error {
	return m.messageHub.Push(msg)
}

func (m *Manager) BroadcastTo(channel string, event string, payload interface{}) error {
	return m.messageHub.Push(&contracts.Message{
		Channel: channel,
		Event:   event,
		Payload: payload,
	})
}

// Handle handles the client
func (m *Manager) Handle(client *Client) {
	defer m.RemoveClient(client.SocketID())

	m.AddClient(client)

	for {
		select {
		case channel := <-client.Subscribe():
			m.broadcaster.Subscribe(channel, client)
		case channel := <-client.Unsubscribe():
			m.broadcaster.Unsubscribe(channel, client)
		case msg := <-client.Receive():
			err := m.messageHub.Push(msg)
			if err != nil {
				// todo: retry or log
			}
		case <-client.Closed():
			return
		}
	}
}

// Start starts the manager
func (m *Manager) Start() {
	go func() {
		for {
			select {
			case msg := <-m.messageHub.Pop():
				m.broadcaster.Send(msg)
			}
		}
	}()
}
