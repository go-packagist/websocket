package websocket

type Subscription struct {
	Channel Channel
	Client  *Client
}
