package contracts

type Subscriber interface {
	Send(payload *Payload)
}
