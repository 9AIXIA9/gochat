package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicMessageDelivered event.Topic = "notification.message.delivered"

var _ event.SpecificEvent = (*MessageDeliveredEvent)(nil)

type MessageDeliveredEvent struct {
	messageID MessageID
	recipient kernel.UserID
	*event.StandardEvent
}

func ToMessageDeliveredEvent(ev event.Event) (*MessageDeliveredEvent, error) {
	e := &MessageDeliveredEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewMessageDeliveredEvent(
	id event.ID,
	messageID MessageID,
	recipient kernel.UserID,
) (*MessageDeliveredEvent, error) {
	e := &MessageDeliveredEvent{
		messageID: messageID,
		recipient: recipient,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(recipient), time.Now().UTC(), TopicMessageDelivered, payload)
	return e, nil
}

func (e *MessageDeliveredEvent) Marshal() ([]byte, error) {
	type Alias struct {
		MessageID MessageID
		Recipient kernel.UserID
	}
	return json.Marshal(Alias{
		MessageID: e.messageID,
		Recipient: e.recipient,
	})
}

func (e *MessageDeliveredEvent) Unmarshal(data []byte) error {
	type Alias struct {
		MessageID MessageID
		Recipient kernel.UserID
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.messageID = tmp.MessageID
	e.recipient = tmp.Recipient
	return nil
}

func (e *MessageDeliveredEvent) MessageID() MessageID {
	return e.messageID
}

func (e *MessageDeliveredEvent) Recipient() kernel.UserID {
	return e.recipient
}
