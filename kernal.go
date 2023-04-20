package websocket

import "github.com/go-fires/websocket/contracts"

type Kernel struct {
	Broadcaster contracts.Broadcaster
}

func NewKernel(broadcaster contracts.Broadcaster) *Kernel {
	return &Kernel{
		Broadcaster: broadcaster,
	}
}

func (k *Kernel) Broadcast(broadcast contracts.Broadcast) {
	k.Broadcaster.Broadcast(broadcast)
}
