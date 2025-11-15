package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicWelcomeEmailNotificationRequested event.Topic = "notification.welcome_email_notification.requested"

var _ event.SpecificEvent = (*WelcomeEmailRequestedNotificationEvent)(nil)

type WelcomeEmailRequestedNotificationEvent struct {
	email  kernel.Email
	number UserNumber
	*event.StandardEvent
}

func ToWelcomeEmailNotificationRequestedEvent(ev event.Event) (*WelcomeEmailRequestedNotificationEvent, error) {
	e := &WelcomeEmailRequestedNotificationEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewWelcomeEmailNotificationRequestedEvent(
	id event.ID,
	userID kernel.UserID,
	email kernel.Email,
	number UserNumber,
) (*WelcomeEmailRequestedNotificationEvent, error) {
	e := &WelcomeEmailRequestedNotificationEvent{
		email:  email,
		number: number,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(userID), time.Now().UTC(), TopicWelcomeEmailNotificationRequested, payload)
	return e, nil
}

func (e *WelcomeEmailRequestedNotificationEvent) Marshal() ([]byte, error) {
	type Alias struct {
		Email  kernel.Email
		Number UserNumber
	}
	return json.Marshal(&Alias{
		Email:  e.email,
		Number: e.number,
	})
}

func (e *WelcomeEmailRequestedNotificationEvent) Unmarshal(data []byte) error {
	type Alias struct {
		Email  kernel.Email
		Number UserNumber
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.email = tmp.Email
	e.number = tmp.Number
	return nil
}

func (e *WelcomeEmailRequestedNotificationEvent) Email() kernel.Email {
	return e.email
}

func (e *WelcomeEmailRequestedNotificationEvent) Number() UserNumber {
	return e.number
}
