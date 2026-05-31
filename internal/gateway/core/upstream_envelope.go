package core

import (
	"encoding/json"
	"gochat/internal/shared/command"
	"gochat/internal/shared/kernel"
)

// UpstreamEnvelope 定义了由客户端发来的标准信封结构
type UpstreamEnvelope struct {
	ClientMessageID kernel.MessageID `json:"client_message_id"`
	Action          command.Action   `json:"action"`
	Payload         json.RawMessage  `json:"payload"`
}
