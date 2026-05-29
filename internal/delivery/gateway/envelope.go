package gateway

import (
	"encoding/json"
	"gochat/internal/shared/event"
)

// WSEnvelope 定义了由客户端发来的标准信封结构
type WSEnvelope struct {
	ClientMessageID string          `json:"client_message_id"`
	Action          event.Topic     `json:"action"` // 它将作为 Kafka 的 Topic（或路由键）
	Payload         json.RawMessage `json:"payload"`
}
