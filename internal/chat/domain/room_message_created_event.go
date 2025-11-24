package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicRoomMessageCreated event.Topic = "chat.room_message.created"

var _ event.SpecificEvent = (*RoomMessageCreatedEvent)(nil)

type RoomMessageCreatedEvent struct {
	*event.StandardEvent
}

func ToRoomMessageCreatedEvent(ev event.Event) (*RoomMessageCreatedEvent, error) {
	e := &RoomMessageCreatedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewRoomMessageCreatedEvent(
	messageID kernel.MessageID,
	generator event.IDGenerator,
) (*RoomMessageCreatedEvent, error) {
	e := &RoomMessageCreatedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(messageID), TopicRoomMessageCreated, payload, generator)
	return e, nil
}

func (e *RoomMessageCreatedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *RoomMessageCreatedEvent) Unmarshal([]byte) error {
	return nil
}
