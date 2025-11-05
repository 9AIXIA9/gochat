package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

type PrivateMessage struct {
	id MessageID
	*MessageInformation
	*RecipientMessageState
}

func NewPrivateMessage(
	id MessageID,
	sender kernel.UserID,
	content string,
	sentAt time.Time,
	state State,
	recipient kernel.UserID,
) *PrivateMessage {
	return &PrivateMessage{
		id:                    id,
		MessageInformation:    NewMessageInformation(sender, content, sentAt),
		RecipientMessageState: NewRecipientMessageState(recipient, state),
	}
}

func (m *PrivateMessage) ID() MessageID {
	return m.id
}
