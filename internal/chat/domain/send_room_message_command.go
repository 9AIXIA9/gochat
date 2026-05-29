package domain

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicSendRoomMessageCommand event.Topic = "send_room_message_command"

var _ event.SpecificEvent = (*SendRoomMessageCommand)(nil)

type SendRoomMessageCommand struct {
	roomID  kernel.RoomID
	content string
	*event.StandardEvent
}

func ToSendRoomMessageCommand(ev event.Event) (*SendRoomMessageCommand, error) {
	if ev.Topic() != TopicSendRoomMessageCommand {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &SendRoomMessageCommand{StandardEvent: event.LoadStandardEventFromEvent(ev)}
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
		RoomID  kernel.RoomID
		Content string
	}
	return json.Marshal(Alias{
		RoomID:  e.roomID,
		Content: e.content,
	})
}

func (e *SendRoomMessageCommand) Unmarshal(data []byte) error {
	type Alias struct {
		RoomID  kernel.RoomID
		Content string
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.roomID = tmp.RoomID
	e.content = tmp.Content
	return nil
}
