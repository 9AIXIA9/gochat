package domain

import (
	"encoding/json"
	"gochat/internal/shared/command"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

const ActionSendRoomMessageCommand command.Action = "send_room_message"

var _ command.SpecificCommand = (*SendRoomMessageCommand)(nil)

type SendRoomMessageCommand struct {
	roomID  kernel.RoomID
	content string
	*command.StandardCommand
}

func ToSendRoomMessageCommand(ev command.Command) (*SendRoomMessageCommand, error) {
	if ev.Action() != ActionSendRoomMessageCommand {
		return nil, myErrors.ErrWrongCommandAction
	}
	e := &SendRoomMessageCommand{StandardCommand: command.LoadStandardCommandFromCommand(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *SendRoomMessageCommand) Content() string {
	return e.content
}

func (e *SendRoomMessageCommand) RoomID() kernel.RoomID {
	return e.roomID
}

func (e *SendRoomMessageCommand) Marshal() ([]byte, error) {
	type Alias struct {
		RoomID  kernel.RoomID `json:"room_id"`
		Content string        `json:"content"`
	}
	return json.Marshal(Alias{
		RoomID:  e.roomID,
		Content: e.content,
	})
}

func (e *SendRoomMessageCommand) Unmarshal(data []byte) error {
	type Alias struct {
		RoomID  kernel.RoomID `json:"room_id"`
		Content string        `json:"content"`
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.roomID = tmp.RoomID
	e.content = tmp.Content
	return nil
}
