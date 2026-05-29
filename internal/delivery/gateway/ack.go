package gateway

import (
	"encoding/json"
)

type AckType string

const (
	AckReceived AckType = "received"
	AckError    AckType = "error"
)

// Ack 定义了发给客户端的快速回执
type Ack struct {
	ClientMsgID string  `json:"clientMsgID"`
	AckType     AckType `json:"ackType"` // "received" or "error"
}

func NewAckReceived(clientMsgID string) ([]byte, error) {
	ack := &Ack{
		ClientMsgID: clientMsgID,
		AckType:     AckReceived,
	}
	return json.Marshal(ack)
}

func NewAckError(clientMsgID string) ([]byte, error) {
	ack := &Ack{
		ClientMsgID: clientMsgID,
		AckType:     AckError,
	}
	return json.Marshal(ack)
}
