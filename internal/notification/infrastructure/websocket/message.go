package websocket

import "encoding/json"

type MessageType string

const (
	PingType MessageType = "ping"
	PongType             = "pong"
	DataType             = "data"
)

type Message struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}
