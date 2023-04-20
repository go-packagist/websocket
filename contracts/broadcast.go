package contracts

type Broadcast interface {
	Channels() []Channeler
	Event() string
	Payload() interface{}
}
