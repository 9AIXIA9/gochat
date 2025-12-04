package domain

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicSystemMessageNotificationRequested event.Topic = "notification.system_message_notification.requested"

var _ event.SpecificEvent = (*SystemMessageNotificationRequestedEvent)(nil)

type SystemMessageNotificationRequestedEvent struct {
	content string
	*event.StandardEvent
}

func ToSystemMessageNotificationRequestedEvent(ev event.Event) (*SystemMessageNotificationRequestedEvent, error) {
	if ev.Topic() != TopicSystemMessageNotificationRequested {
		return nil, myErrors.ErrWrongEventType
	}
	e := &SystemMessageNotificationRequestedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewSystemMessageNotificationRequestedEvent(
	recipientID kernel.UserID,
	content string,
	generator event.IDGenerator,
) (*SystemMessageNotificationRequestedEvent, error) {
	e := &SystemMessageNotificationRequestedEvent{
		content: content,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(
		kernel.ID(recipientID),
		TopicSystemMessageNotificationRequested,
		payload,
		generator,
	)
	return e, nil
}

func (e *SystemMessageNotificationRequestedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		Content string
	}
	return json.Marshal(Alias{
		Content: e.content,
	})
}

func (e *SystemMessageNotificationRequestedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		Content string
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.content = tmp.Content
	return nil
}

func (e *SystemMessageNotificationRequestedEvent) Content() string {
	return e.content
}
