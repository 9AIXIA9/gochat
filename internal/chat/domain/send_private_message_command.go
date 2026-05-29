package domain

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicSendPrivateMessageCommand event.Topic = "chat.send_private_message_command"

var _ event.SpecificEvent = (*SendPrivateMessageCommand)(nil)

type SendPrivateMessageCommand struct {
	recipientID kernel.UserID
	content     string
	*event.StandardEvent
}

func ToSendPrivateMessageCommand(ev event.Event) (*SendPrivateMessageCommand, error) {
	if ev.Topic() != TopicSendPrivateMessageCommand {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &SendPrivateMessageCommand{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *SendPrivateMessageCommand) Content() string {
	return e.content
}

func (e *SendPrivateMessageCommand) RecipientID() kernel.UserID {
	return e.recipientID
}

func (e *SendPrivateMessageCommand) Marshal() ([]byte, error) {
	type Alias struct {
		RecipientID kernel.UserID
		Content     string
	}
	return json.Marshal(Alias{
		RecipientID: e.recipientID,
		Content:     e.content,
	})
}

func (e *SendPrivateMessageCommand) Unmarshal(data []byte) error {
	type Alias struct {
		RecipientID kernel.UserID
		Content     string
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.recipientID = tmp.RecipientID
	e.content = tmp.Content
	return nil
}
