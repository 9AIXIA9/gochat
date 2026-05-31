package domain

import (
	"encoding/json"
	"gochat/internal/shared/command"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

const ActionSendPrivateMessageCommand command.Action = "send_private_message"

var _ command.SpecificCommand = (*SendPrivateMessageCommand)(nil)

type SendPrivateMessageCommand struct {
	recipientID kernel.UserID
	content     string
	*command.StandardCommand
}

func ToSendPrivateMessageCommand(com command.Command) (*SendPrivateMessageCommand, error) {
	if com.Action() != ActionSendPrivateMessageCommand {
		return nil, myErrors.ErrWrongCommandAction
	}
	e := &SendPrivateMessageCommand{StandardCommand: command.LoadStandardCommandFromCommand(com)}
	if len(com.Payload()) > 0 {
		if err := e.Unmarshal(com.Payload()); err != nil {
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
		RecipientID kernel.UserID `json:"recipient_id"`
		Content     string        `json:"content"`
	}
	return json.Marshal(Alias{
		RecipientID: e.recipientID,
		Content:     e.content,
	})
}

func (e *SendPrivateMessageCommand) Unmarshal(data []byte) error {
	type Alias struct {
		RecipientID kernel.UserID `json:"recipient_id"`
		Content     string        `json:"content"`
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.recipientID = tmp.RecipientID
	e.content = tmp.Content
	return nil
}
