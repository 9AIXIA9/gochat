package domain

import (
	"context"
	"gochat/internal/shared/event"
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

	eventManager *event.Manager
}

func LoadPrivateMessage(
	id kernel.MessageID,
	senderID kernel.UserID,
	recipientID kernel.UserID,
	content string,
	state MessageState,
	sentAt time.Time,
) *PrivateMessage {
	return &PrivateMessage{
		id:           id,
		senderID:     senderID,
		recipientID:  recipientID,
		state:        state,
		content:      content,
		sentAt:       sentAt,
		eventManager: event.NewEventManager(),
	}
}

func CreatePrivateMessage(
	ctx context.Context,
	recipientID kernel.UserID,
	senderID kernel.UserID,
	content string,
	messageIDGenerator kernel.MessageIDGenerator,
	notifier PrivateMessageNotifier,
	exister FriendshipExisterByUserID,
) (*PrivateMessage, error) {
	if len(content) == 0 {
		return nil, ErrEmptyMessageContent
	}

	if senderID != recipientID {
		exist, err := exister.ExistByUserID(ctx, senderID, recipientID)
		if err != nil {
			return nil, err
		}

		if !exist {
			return nil, ErrNotFriends
		}
	}

	message := &PrivateMessage{
		id:           messageIDGenerator.Generate(),
		senderID:     senderID,
		recipientID:  recipientID,
		state:        MessageStateUndelivered,
		content:      content,
		sentAt:       time.Now().UTC(),
		eventManager: event.NewEventManager(),
	}

	if err := notifier.Notify(message); err != nil {
		return message, nil
	}

	return message, nil
}

func (m *PrivateMessage) Deliver(
	notifier PrivateMessageNotifier,
) error {
	if m.state != MessageStateUndelivered {
		return nil
	}

	if err := notifier.Notify(m); err != nil {
		return err
	}

	m.state = MessageStateDelivered

	return nil
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

func (m *PrivateMessage) Content() string {
	return m.content
}

func (m *PrivateMessage) State() MessageState {
	return m.state
}

func (m *PrivateMessage) SentAt() time.Time {
	return m.sentAt
}

func (m *PrivateMessage) GetEvents() []event.Event {
	return m.eventManager.GetEvents()
}
