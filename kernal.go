package websocket

import (
	"github.com/go-fires/websocket/contracts"
	"log"
)

type Kernel struct {
	broadcaster contracts.Broadcaster

	messageHub contracts.MessageHub
}

func NewKernel(broadcaster contracts.Broadcaster, messageHub contracts.MessageHub) *Kernel {
	k := &Kernel{
		broadcaster: broadcaster,
		messageHub:  messageHub,
	}

	go k.send()

	return k
}

func (k *Kernel) send() {
	for {
		select {
		case msg := <-k.messageHub.Pop():
			k.broadcaster.Send(msg)
		}
	}
}

func (k *Kernel) Broadcast(broadcast contracts.Broadcast) {
	for _, channel := range broadcast.Channels() {
		err := k.messageHub.Push(&contracts.Message{
			Channel: channel.Name(),
			Event:   broadcast.Event(),
			Payload: broadcast.Payload(),
		})
		if err != nil {
			log.Printf("websocket: failed to push message to message hub: %v", err)
			continue
		}
	}
}
