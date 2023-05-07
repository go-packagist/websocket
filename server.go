package websocket

import (
	"time"
)

type Server struct {
	clients map[string]*Client
	handler Handler
	broker  Broker
}

func NewServer(handler Handler, broker Broker) *Server {
	return &Server{
		clients: make(map[string]*Client),
		handler: handler,
		broker:  broker,
	}
}

func (s *Server) Start() {
	defer s.stop()

	for {
		select {
		case message := <-s.broker.Consume():
			s.Emit(message)
		}
	}
}

// Broadcast  a message to all subscriber clients in cluster nodes
func (s *Server) Broadcast(message *Message) {
	s.broker.Publish() <- message
}

// Emit a message to all subscriber clients in current node
func (s *Server) Emit(message *Message) {
	if message.Payload != nil && message.Payload.Channel != "" {
		for _, client := range s.clients {
			if client.In(message.Payload.Channel) {
				client.Write() <- message
			}
		}
	}
}

func (s *Server) stop() {
	for _, client := range s.clients {
		client.Close()
	}
}
func (s *Server) Handle(client *Client) {
	s.clients[client.ID()] = client
	defer func() {
		delete(s.clients, client.ID())
	}()

	for {
		select {
		case message := <-client.Read():
			s.handler.Message(client, message)
		case <-client.Closed():
			s.handler.Closed(client)
			return
		case err := <-client.Error():
			s.handler.Error(client, err)
		case <-time.Tick(25 * time.Second):
			client.Ping()
		}
	}
}
