package core

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type CommandResult string

func (r CommandResult) String() string {
	return string(r)
}

const (
	CommandResultSuccess CommandResult = "success"
	CommandResultFailed  CommandResult = "failed"
)

// DownstreamEnvelop 定义了由服务器发往客户端的标准信封结构
type DownstreamEnvelop struct {
	ClientMessageID kernel.MessageID `json:"client_message_id"`
	Action          event.Topic      `json:"action"`
	Result          CommandResult    `json:"result"`
	Payload         json.RawMessage  `json:"payload"`
	ErrorMessage    string           `json:"error_message,omitempty"`
}

func NewSuccessEnvelop(clientMessageID kernel.MessageID, action event.Topic, payload []byte) *DownstreamEnvelop {
	return &DownstreamEnvelop{
		ClientMessageID: clientMessageID,
		Action:          action,
		Result:          CommandResultSuccess,
		Payload:         payload,
	}
}

func NewFailedEnvelop(clientMessageID kernel.MessageID, action event.Topic, errorMessage string) *DownstreamEnvelop {
	return &DownstreamEnvelop{
		ClientMessageID: clientMessageID,
		Action:          action,
		Result:          CommandResultFailed,
		ErrorMessage:    errorMessage,
	}
}
