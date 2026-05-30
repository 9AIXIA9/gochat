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
	Action          event.Topic      `json:"action"`
	Payload         json.RawMessage  `json:"payload"`
	ClientMessageID kernel.MessageID `json:"client_message_id,omitempty"`
	Result          CommandResult    `json:"result,omitempty"`
	ErrorMessage    string           `json:"error_message,omitempty"`
}

func NewActiveDownstreamEnvelop(action event.Topic, payload json.RawMessage) *DownstreamEnvelop {
	return &DownstreamEnvelop{
		Action:  action,
		Payload: payload,
	}
}

func NewActionSucceededDownstreamEnvelop(clientMessageID kernel.MessageID, action event.Topic, payload []byte) *DownstreamEnvelop {
	return &DownstreamEnvelop{
		ClientMessageID: clientMessageID,
		Action:          action,
		Result:          CommandResultSuccess,
		Payload:         payload,
	}
}

func NewActionFailedDownstreamEnvelop(clientMessageID kernel.MessageID, action event.Topic, errorMessage string) *DownstreamEnvelop {
	return &DownstreamEnvelop{
		ClientMessageID: clientMessageID,
		Action:          action,
		Result:          CommandResultFailed,
		ErrorMessage:    errorMessage,
	}
}
