package websocket

type Message struct {
	Topic Topic `json:"topic"`
	Body  any   `json:"body,omitempty"`
}
