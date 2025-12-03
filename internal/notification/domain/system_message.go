package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

type SystemMessage struct {
	id          kernel.MessageID
	recipientID kernel.UserID
	state       MessageState
	content     []byte
	sentAt      time.Time
}

func LoadSystemMessage(
	id kernel.MessageID,
	recipientID kernel.UserID,
	state MessageState,
	content []byte,
	sentAt time.Time,
) *SystemMessage {
	return &SystemMessage{
		id:          id,
		recipientID: recipientID,
		state:       state,
		content:     content,
		sentAt:      sentAt,
	}
}

func CreateSystemMessage(
	recipientID kernel.UserID,
	content []byte,
	idGenerator MessageIDGenerator,
	notifier SystemMessageNotifier,
) (*SystemMessage, error) {
	if len(content) == 0 {
		return nil, ErrEmptyContent
	}

	message := &SystemMessage{
		id:          idGenerator.Generate(),
		recipientID: recipientID,
		state:       MessageStateUndelivered,
		content:     content,
		sentAt:      time.Now().UTC(),
	}

	//尝试投递消息，投递失败不影响消息创建
	_ = message.Deliver(notifier)

	return message, nil
}

func (m *SystemMessage) Deliver(
	notifier SystemMessageNotifier,
) error {
	if m.state != MessageStateUndelivered {
		return ErrNotUndelivered
	}

	if err := notifier.Notify(m); err != nil {
		return err
	}

	m.state = MessageStateDelivered

	return nil
}

func (m *SystemMessage) ID() kernel.MessageID {
	return m.id
}

func (m *SystemMessage) RecipientID() kernel.UserID {
	return m.recipientID
}

func (m *SystemMessage) State() MessageState {
	return m.state
}

func (m *SystemMessage) Content() []byte {
	return m.content
}

func (m *SystemMessage) SentAt() time.Time {
	return m.sentAt
}
