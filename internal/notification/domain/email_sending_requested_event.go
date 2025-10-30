package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicEmailSendingRequested = "notification.email.sending.requested"

var _ event.SpecificEvent = (*EmailSendingRequestedEvent)(nil)

type EmailSendingRequestedEvent struct {
	email   kernel.Email
	theme   NoticeTheme
	title   string
	content string
	*event.StandardEvent
}

func ToEmailSendingRequestedEvent(ev event.Event) (*EmailSendingRequestedEvent, error) {
	e := &EmailSendingRequestedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewEmailSendingRequestedEvent(id event.ID, userID kernel.UserID, email kernel.Email, theme NoticeTheme, title string, content string) (*EmailSendingRequestedEvent, error) {
	e := &EmailSendingRequestedEvent{
		email:   email,
		theme:   theme,
		title:   title,
		content: content,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(userID), time.Now().UTC(), TopicEmailSendingRequested, payload)
	return e, nil
}

func (e *EmailSendingRequestedEvent) Email() kernel.Email {
	return e.email
}

func (e *EmailSendingRequestedEvent) Theme() NoticeTheme {
	return e.theme
}
func (e *EmailSendingRequestedEvent) Title() string {
	return e.title
}
func (e *EmailSendingRequestedEvent) Content() string {
	return e.content
}

func (e *EmailSendingRequestedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		Email   kernel.Email
		Theme   NoticeTheme
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

func (e *EmailSendingRequestedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		Email   kernel.Email
		Theme   NoticeTheme
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
