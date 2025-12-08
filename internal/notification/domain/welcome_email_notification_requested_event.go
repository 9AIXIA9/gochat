package domain

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicWelcomeEmailNotificationRequested event.Topic = "notification.welcome_email_notification.requested"

var _ event.SpecificEvent = (*WelcomeEmailRequestedNotificationEvent)(nil)

type WelcomeEmailRequestedNotificationEvent struct {
	email  kernel.Email
	number kernel.UserNumber
	*event.StandardEvent
}

func ToWelcomeEmailNotificationRequestedEvent(ev event.Event) (*WelcomeEmailRequestedNotificationEvent, error) {
	if ev.Topic() != TopicWelcomeEmailNotificationRequested {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &WelcomeEmailRequestedNotificationEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewWelcomeEmailNotificationRequestedEvent(
	userID kernel.UserID,
	email kernel.Email,
	number kernel.UserNumber,
	generator event.IDGenerator,
) (*WelcomeEmailRequestedNotificationEvent, error) {
	e := &WelcomeEmailRequestedNotificationEvent{
		email:  email,
		number: number,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(userID), TopicWelcomeEmailNotificationRequested, payload, generator)
	return e, nil
}

func (e *WelcomeEmailRequestedNotificationEvent) Marshal() ([]byte, error) {
	type Alias struct {
		Email  kernel.Email
		Number kernel.UserNumber
	}
	return json.Marshal(&Alias{
		Email:  e.email,
		Number: e.number,
	})
}

func (e *WelcomeEmailRequestedNotificationEvent) Unmarshal(data []byte) error {
	type Alias struct {
		Email  kernel.Email
		Number kernel.UserNumber
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

func (e *WelcomeEmailRequestedNotificationEvent) Number() kernel.UserNumber {
	return e.number
}
