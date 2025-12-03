package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

const (
	TopicFriendshipCreated    SystemMessageTopic = "friendship_created"
	TopicFriendRequestCreated SystemMessageTopic = "friend_request_created"
)

type SystemMessageTopic string

func (t SystemMessageTopic) String() string {
	return string(t)
}

type SystemMessage struct {
	id          kernel.MessageID
	topic       SystemMessageTopic
	recipientID kernel.UserID
	state       MessageState
	content     []byte
	sentAt      time.Time
}

func LoadSystemMessage(
	id kernel.MessageID,
	topic SystemMessageTopic,
	recipientID kernel.UserID,
	state MessageState,
	content []byte,
	sentAt time.Time,
) *SystemMessage {
	return &SystemMessage{
		id:          id,
		topic:       topic,
		recipientID: recipientID,
		state:       state,
		content:     content,
		sentAt:      sentAt,
	}
}

func CreateSystemMessage(
	topic SystemMessageTopic,
	recipientID kernel.UserID,
	content []byte,
	idGenerator kernel.MessageIDGenerator,
	notifier SystemMessageNotifier,
) (*SystemMessage, error) {
	if len(content) == 0 {
		return nil, ErrEmptyContent
	}

	message := &SystemMessage{
		id:          idGenerator.Generate(),
		topic:       topic,
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

	m.state = MessageStateDelivered

	if err := notifier.Notify(m); err != nil {
		m.state = MessageStateUndelivered
		return err
	}

	return nil
}

func (m *SystemMessage) ID() kernel.MessageID {
	return m.id
}

func (m *SystemMessage) Topic() SystemMessageTopic {
	return m.topic
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
