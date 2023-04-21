package contracts

type Handler interface {
	OnOpen(client Client)
	OnClose(client Client)
	OnMessage(client Client, message *Message)
	OnError(client Client, err error)
}
