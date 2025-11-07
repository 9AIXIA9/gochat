package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicPrivateMessageReceived event.Topic = "chat.private_message.received"

var _ event.SpecificEvent = (*PrivateMessageReceivedEvent)(nil)

type PrivateMessageReceivedEvent struct {
	messageID MessageID
	*event.StandardEvent
}

func ToPrivateMessageReceivedEvent(ev event.Event) (*PrivateMessageReceivedEvent, error) {
	e := &PrivateMessageReceivedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewPrivateMessageReceivedEvent(
	id event.ID,
	messageID MessageID,
	recipient kernel.UserID,
) (*PrivateMessageReceivedEvent, error) {
	e := &PrivateMessageReceivedEvent{
		messageID: messageID,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(recipient), time.Now().UTC(), TopicPrivateMessageReceived, payload)
	return e, nil
}

func (e *PrivateMessageReceivedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		MessageID MessageID
	}
	return json.Marshal(Alias{
		MessageID: e.messageID,
	})
}

func (e *PrivateMessageReceivedEvent) Unmarshal(data []byte) error {
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

func (e *PrivateMessageReceivedEvent) MessageID() MessageID {
	return e.messageID
}
