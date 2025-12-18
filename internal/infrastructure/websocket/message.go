package websocket

import "encoding/json"

// Message 表示 WebSocket 通用消息
type Message struct {
	Topic Topic `json:"topic"`
	Body  any   `json:"body,omitempty"`
}

// Request 表示 WebSocket 请求。
// 使用 json.RawMessage 延迟解码，按 Topic 路由后再做具体反序列化。
type Request struct {
	Topic   Topic           `json:"topic" validate:"required"`
	Payload json.RawMessage `json:"payload,omitempty"`
}
