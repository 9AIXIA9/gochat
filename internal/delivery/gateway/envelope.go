package gateway

import (
	"encoding/json"
	"gochat/internal/shared/event"
)

// Envelope 定义了由客户端发来的标准信封结构
type Envelope struct {
	ClientMessageID string          `json:"client_message_id"`
	Action          event.Topic     `json:"action"`
	Payload         json.RawMessage `json:"payload"`
}
