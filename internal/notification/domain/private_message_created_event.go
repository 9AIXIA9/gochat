package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicPrivateMessageCreated event.Topic = "notification.private_message.created"

var _ event.SpecificEvent = (*PrivateMessageCreatedEvent)(nil)

type PrivateMessageCreatedEvent struct {
	messageID MessageID
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
	messageID MessageID,
	recipient kernel.UserID,
) (*PrivateMessageCreatedEvent, error) {
	e := &PrivateMessageCreatedEvent{
		messageID: messageID,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(recipient), time.Now().UTC(), TopicPrivateMessageCreated, payload)
	return e, nil
}

func (e *PrivateMessageCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		MessageID MessageID
	}
	return json.Marshal(Alias{
		MessageID: e.messageID,
	})
}

func (e *PrivateMessageCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		MessageID MessageID
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.messageID = tmp.MessageID
	return nil
}

func (e *PrivateMessageCreatedEvent) MessageID() MessageID {
	return e.messageID
}
