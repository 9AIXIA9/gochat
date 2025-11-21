package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

type PrivateMessage struct {
	id          kernel.MessageID
	senderID    kernel.UserID
	recipientID kernel.UserID
	state       MessageState
	content     string
	sentAt      time.Time
}

func LoadPrivateMessage(
	id kernel.MessageID,
	senderID kernel.UserID,
	recipientID kernel.UserID,
	state MessageState,
	content string,
	sentAt time.Time,
) *PrivateMessage {
	return &PrivateMessage{
		id:          id,
		senderID:    senderID,
		recipientID: recipientID,
		state:       state,
		content:     content,
		sentAt:      sentAt,
	}
}

func ReceivePrivateMessage(
	id kernel.MessageID,
	senderID kernel.UserID,
	recipientID kernel.UserID,
	content string,
	sentAt time.Time,
) *PrivateMessage {
	return &PrivateMessage{
		id:          id,
		senderID:    senderID,
		recipientID: recipientID,
		state:       MessageStateReceived,
		content:     content,
		sentAt:      sentAt,
	}
}

func (m *PrivateMessage) Deliver(
	notifier PrivateMessageNotifier,
) error {
	if err := notifier.NotifyPrivateMessage(m); err != nil {
		return err
	}

	m.state = MessageStateDelivered
	return nil
}

func (m *PrivateMessage) Read() {
	m.state = MessageStateRead
}

func (m *PrivateMessage) ID() kernel.MessageID {
	return m.id
}

func (m *PrivateMessage) SenderID() kernel.UserID {
	return m.senderID
}

func (m *PrivateMessage) RecipientID() kernel.UserID {
	return m.recipientID
}

func (m *PrivateMessage) State() MessageState {
	return m.state
}

func (m *PrivateMessage) Content() string {
	return m.content
}

func (m *PrivateMessage) SentAt() time.Time {
	return m.sentAt
}
