package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicEmailNotificationRequested = "notification.email_notification.requested"

var _ event.SpecificEvent = (*EmailNotificationRequestedEvent)(nil)

type EmailNotificationRequestedEvent struct {
	email   kernel.Email
	theme   MailTheme
	title   string
	content string
	*event.StandardEvent
}

func ToEmailNotificationRequestedEvent(ev event.Event) (*EmailNotificationRequestedEvent, error) {
	e := &EmailNotificationRequestedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewEmailNotificationRequestedEvent(id event.ID, userID kernel.UserID, email kernel.Email, theme MailTheme, title string, content string) (*EmailNotificationRequestedEvent, error) {
	e := &EmailNotificationRequestedEvent{
		email:   email,
		theme:   theme,
		title:   title,
		content: content,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(userID), time.Now().UTC(), TopicEmailNotificationRequested, payload)
	return e, nil
}

func (e *EmailNotificationRequestedEvent) Email() kernel.Email {
	return e.email
}

func (e *EmailNotificationRequestedEvent) Theme() MailTheme {
	return e.theme
}
func (e *EmailNotificationRequestedEvent) Title() string {
	return e.title
}
func (e *EmailNotificationRequestedEvent) Content() string {
	return e.content
}

func (e *EmailNotificationRequestedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		Email   kernel.Email
		Theme   MailTheme
		Title   string
		Content string
	}
	return json.Marshal(Alias{
		Email:   e.email,
		Theme:   e.theme,
		Title:   e.title,
		Content: e.content,
	})
}

func (e *EmailNotificationRequestedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		Email   kernel.Email
		Theme   MailTheme
		Title   string
		Content string
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.email = tmp.Email
	e.theme = tmp.Theme
	e.title = tmp.Title
	e.content = tmp.Content
	return nil
}
