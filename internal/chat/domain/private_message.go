package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

type PrivateMessage struct {
	*MessageInformation
	recipient kernel.UserID
	state     MessageState
}

func NewPrivateMessage(
	id MessageID,
	sender kernel.UserID,
	content string,
	sentAt time.Time,
	state MessageState,
	recipient kernel.UserID,
) *PrivateMessage {
	return &PrivateMessage{
		MessageInformation: NewMessageInformation(id, PrivateType, sender, content, sentAt),
		recipient:          recipient,
		state:              state,
	}
}

func (m *PrivateMessage) Recipient() kernel.UserID {
	return m.recipient
}

func (m *PrivateMessage) State() MessageState {
	return m.state
}
