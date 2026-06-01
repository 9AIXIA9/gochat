package core

import (
	"encoding/json"
	"gochat/internal/shared/command"
	"gochat/internal/shared/kernel"
)

// Envelope 标准信封结构
type Envelope struct {
	ClientMessageID kernel.MessageID `json:"client_message_id,omitempty"`
	Action          command.Action   `json:"action"`
	Payload         json.RawMessage  `json:"payload"`
}

func NewEnvelop(
	action command.Action,
	payload json.RawMessage,
) *Envelope {
	return &Envelope{
		Action:  action,
		Payload: payload,
	}
}

func NewEnvelopWithClientMessageID(
	action command.Action,
	payload json.RawMessage,
	clientMessageID kernel.MessageID,
) *Envelope {
	return &Envelope{
		ClientMessageID: clientMessageID,
		Action:          action,
		Payload:         payload,
	}
}
