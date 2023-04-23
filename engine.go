package websocket

import (
	"fmt"
	"sync"
	"time"
)

type Engine struct {
	emitter *Emitter
	hub     Hub
	handler Handler

	mu sync.Mutex
}

func NewEngine(hub Hub, handler Handler) *Engine {
	return &Engine{
		emitter: NewEmitter(),
		hub:     hub,
		handler: handler,
	}
}

func (e *Engine) Start() {
	go e.monitor()

	for {
		select {
		case msg := <-e.hub.Read():
			e.emitter.Emit(msg)
		}
	}
}

func (e *Engine) Broadcast(msg *Message) {
	e.hub.Recv() <- msg
}

// Handle handle a socket
func (e *Engine) Handle(socketer *Socketer) {
	defer e.emitter.Unsubscribe(socketer)
	e.emitter.Subscribe(socketer)

	e.handler.Opened(socketer)
	defer e.handler.Closed(socketer)

	for {
		select {
		case <-socketer.Closed():
			return
		case err := <-socketer.Error():
			e.handler.Error(socketer, err)
		case msg := <-socketer.Read():
			e.handler.Message(socketer, msg)
		}
	}
}

func (e *Engine) monitor() {
	for {
		select {
		case <-time.Tick(5 * time.Second):
			fmt.Printf("socketer count: %d\n", e.emitter.Count())
		}
	}
}
