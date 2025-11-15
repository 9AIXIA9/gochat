package websocket

import "encoding/json"

type MessageType string

const (
	PingType     MessageType = "ping"
	PongType                 = "pong"
	RequestType              = "request"
	ResponseType             = "response"
)

type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}
