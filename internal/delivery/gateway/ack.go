package gateway

import (
	"encoding/json"
	"gochat/internal/shared/kernel"
)

type AckType string

const (
	AckReceived AckType = "received"
	AckError    AckType = "error"
)

// Ack 定义了发给客户端的快速回执
type Ack struct {
	ClientMessageID kernel.MessageID `json:"client_message_id"`
	AckType         AckType          `json:"ack_type"`        // "received" or "error"
	Error           string           `json:"error,omitempty"` // 详细的错误原因
}

func NewAckReceived(clientMessageID kernel.MessageID) ([]byte, error) {
	ack := &Ack{
		ClientMessageID: clientMessageID,
		AckType:         AckReceived,
	}
	return json.Marshal(ack)
}

func NewAckError(clientMessageID kernel.MessageID, reason error) ([]byte, error) {
	ack := &Ack{
		ClientMessageID: clientMessageID,
		AckType:         AckError,
		Error:           reason.Error(),
	}
	return json.Marshal(ack)
}
