package domain

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicFriendRequestCreatedNotificationRequested event.Topic = "notification.friend_request_created_notification.requested"

var _ event.SpecificEvent = (*FriendRequestCreatedNotificationRequestedEvent)(nil)

type FriendRequestCreatedNotificationRequestedEvent struct {
	from    kernel.UserID
	to      kernel.UserID
	sentAt  time.Time
	content string
	*event.StandardEvent
}

func ToFriendRequestCreatedNotificationRequestedEvent(ev event.Event) (*FriendRequestCreatedNotificationRequestedEvent, error) {
	if ev.Topic() != TopicFriendRequestCreatedNotificationRequested {
		return nil, myErrors.ErrWrongEventType
	}
	e := &FriendRequestCreatedNotificationRequestedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewFriendRequestCreatedNotificationRequestedEvent(
	requestID kernel.OperationID,
	from kernel.UserID,
	to kernel.UserID,
	sentAt time.Time,
	content string,
	generator event.IDGenerator,
) (*FriendRequestCreatedNotificationRequestedEvent, error) {
	e := &FriendRequestCreatedNotificationRequestedEvent{
		from:    from,
		to:      to,
		sentAt:  sentAt,
		content: content,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(
		kernel.ID(requestID),
		TopicFriendRequestCreatedNotificationRequested,
		payload,
		generator,
	)
	return e, nil
}

func (e *FriendRequestCreatedNotificationRequestedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		From    kernel.UserID
		To      kernel.UserID
		SentAt  time.Time
		Content string
	}
	return json.Marshal(Alias{
		From:    e.from,
		To:      e.to,
		SentAt:  e.sentAt,
		Content: e.content,
	})
}

func (e *FriendRequestCreatedNotificationRequestedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		From    kernel.UserID
		To      kernel.UserID
		SentAt  time.Time
		Content string
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.from = tmp.From
	e.to = tmp.To
	e.sentAt = tmp.SentAt
	e.content = tmp.Content
	return nil
}

func (e *FriendRequestCreatedNotificationRequestedEvent) From() kernel.UserID {
	return e.from
}

func (e *FriendRequestCreatedNotificationRequestedEvent) To() kernel.UserID {
	return e.to
}

func (e *FriendRequestCreatedNotificationRequestedEvent) SentAt() time.Time {
	return e.sentAt
}

func (e *FriendRequestCreatedNotificationRequestedEvent) Content() string {
	return e.content
}
