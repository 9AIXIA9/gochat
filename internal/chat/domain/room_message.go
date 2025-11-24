package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

type RoomMessage struct {
	id       kernel.MessageID
	senderID kernel.UserID
	roomID   kernel.RoomID
	content  string
	sentAt   time.Time

	eventManager *event.Manager
}

func LoadRoomMessage(
	id kernel.MessageID,
	senderID kernel.UserID,
	roomID kernel.RoomID,
	content string,
	sentAt time.Time,
) *RoomMessage {
	return &RoomMessage{
		id:           id,
		senderID:     senderID,
		roomID:       roomID,
		content:      content,
		sentAt:       sentAt,
		eventManager: event.NewEventManager(),
	}
}

func CreateRoomMessage(
	room *Room,
	senderID kernel.UserID,
	content string,
	messageIDGenerator MessageIDGenerator,
	eventIDGenerator event.IDGenerator,
) (*RoomMessage, error) {
	if !room.IsMember(senderID) {
		return nil, myErrors.ErrNotBelongTo
	}

	message := &RoomMessage{
		id:           messageIDGenerator.Generate(),
		senderID:     senderID,
		roomID:       room.id,
		content:      content,
		sentAt:       time.Now().UTC(),
		eventManager: event.NewEventManager(),
	}

	ev, err := NewRoomMessageCreatedEvent(
		eventIDGenerator.Generate(),
		message.id,
	)
	if err != nil {
		return nil, err
	}

	message.eventManager.RecordEvent(ev)

	return message, nil
}

func (m *RoomMessage) ID() kernel.MessageID {
	return m.id
}

func (m *RoomMessage) SenderID() kernel.UserID {
	return m.senderID
}

func (m *RoomMessage) RoomID() kernel.RoomID {
	return m.roomID
}

func (m *RoomMessage) Content() string {
	return m.content
}

func (m *RoomMessage) SentAt() time.Time {
	return m.sentAt
}

func (m *RoomMessage) GetEvents() []event.Event {
	return m.eventManager.GetEvents()
}
