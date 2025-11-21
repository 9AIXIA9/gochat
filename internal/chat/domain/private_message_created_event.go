package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

//TODO: 重构事件命名  重构事件处理

const TopicPrivateMessageCreated event.Topic = "chat.private_message.created"

var _ event.SpecificEvent = (*PrivateMessageCreatedEvent)(nil)

type PrivateMessageCreatedEvent struct {
	*event.StandardEvent
}

func ToPrivateMessageCreatedEvent(ev event.Event) (*PrivateMessageCreatedEvent, error) {
	e := &PrivateMessageCreatedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewPrivateMessageCreatedEvent(
	id event.ID,
	messageID kernel.MessageID,
) (*PrivateMessageCreatedEvent, error) {
	e := &PrivateMessageCreatedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(messageID), time.Now().UTC(), TopicPrivateMessageCreated, payload)
	return e, nil
}

func (e *PrivateMessageCreatedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *PrivateMessageCreatedEvent) Unmarshal([]byte) error {
	return nil
}
