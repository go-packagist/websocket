package contracts

type Broadcast interface {
	Channels() []Channel
	Event() string
	Payload() interface{}
}
