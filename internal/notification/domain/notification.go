package domain

import (
	"encoding/json"
	"gochat/internal/shared/kernel"
)

type Notification struct {
	id          kernel.MessageID
	recipientID kernel.UserID
	rawPayload  json.RawMessage
}

func LoadNotification(
	id kernel.MessageID,
	recipientID kernel.UserID,
	rawPayload json.RawMessage,
) *Notification {
	return &Notification{
		id:          id,
		recipientID: recipientID,
		rawPayload:  rawPayload,
	}
}

func CreateNotification(
	id kernel.MessageID,
	recipientID kernel.UserID,
	rawPayload json.RawMessage,
) (*Notification, error) {
	if len(recipientID) == 0 {
		return nil, ErrNotificationNoRecipient
	}
	return &Notification{
		id:          id,
		recipientID: recipientID,
		rawPayload:  rawPayload,
	}, nil
}

func (n *Notification) ID() kernel.MessageID {
	return n.id
}

func (n *Notification) RecipientID() kernel.UserID {
	return n.recipientID
}

func (n *Notification) RawPayload() json.RawMessage {
	return n.rawPayload
}
