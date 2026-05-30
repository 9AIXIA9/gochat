package domain

import (
	"context"
	"encoding/json"
	"gochat/internal/shared/contract"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

type RoomMessage struct {
	id           kernel.MessageID
	senderID     kernel.UserID
	recipientIDs []kernel.UserID
	roomID       kernel.RoomID
	content      string
	sentAt       time.Time

	eventManager *event.Manager
}

func LoadRoomMessage(
	id kernel.MessageID,
	senderID kernel.UserID,
	recipientIDs []kernel.UserID,
	roomID kernel.RoomID,
	content string,
	sentAt time.Time,
) *RoomMessage {
	return &RoomMessage{
		id:           id,
		senderID:     senderID,
		roomID:       roomID,
		content:      content,
		recipientIDs: recipientIDs,
		sentAt:       sentAt,
		eventManager: event.NewEventManager(),
	}
}

func CreateRoomMessage(
	ctx context.Context,
	roomID kernel.RoomID,
	senderID kernel.UserID,
	content string,
	finder RoomshipsFinderByRoomID,
	messageIDGenerator kernel.MessageIDGenerator,
	eventIDGenerator event.IDGenerator,
) (*RoomMessage, error) {
	if len(content) == 0 {
		return nil, ErrEmptyMessageContent
	}

	roomships, err := finder.FindsByRoomID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	message, err := createRoomMessageByRoomships(roomships, roomID, senderID, content, messageIDGenerator)
	if err != nil {
		return nil, err
	}

	rawPayload, err := message.Marshal()
	if err != nil {
		return nil, err
	}

	for _, recipientID := range message.recipientIDs {
		ev, err := contract.NewNotificationCreatedEvent(
			message.id,
			recipientID,
			rawPayload,
			eventIDGenerator,
		)
		if err != nil {
			return nil, err
		}

		message.eventManager.RecordEvent(ev)
	}
	return message, nil
}

func createRoomMessageByRoomships(
	roomships []*Roomship,
	roomID kernel.RoomID,
	senderID kernel.UserID,
	content string,
	messageIDGenerator kernel.MessageIDGenerator,
) (*RoomMessage, error) {
	if len(roomships) == 0 {
		return nil, ErrRoomNotFound
	}

	recipientIDs := make([]kernel.UserID, 0, len(roomships)-1)
	exist := false
	for _, roomship := range roomships {
		if roomship.UserID() == senderID {
			exist = true
			continue
		}
		recipientIDs = append(recipientIDs, roomship.UserID())
	}

	if !exist {
		return nil, ErrNotMember
	}

	if len(recipientIDs) == 0 {
		return &RoomMessage{
			id:           messageIDGenerator.Generate(),
			senderID:     senderID,
			recipientIDs: nil,
			roomID:       roomID,
			content:      content,
			sentAt:       time.Now().UTC(),
			eventManager: event.NewEventManager(),
		}, nil
	}

	return &RoomMessage{
		id:           messageIDGenerator.Generate(),
		senderID:     senderID,
		roomID:       roomID,
		recipientIDs: recipientIDs,
		content:      content,
		sentAt:       time.Now().UTC(),
		eventManager: event.NewEventManager(),
	}, nil
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

func (m *RoomMessage) RecipientIDs() []kernel.UserID {
	return m.recipientIDs
}

func (m *RoomMessage) Marshal() ([]byte, error) {
	type Alias struct {
		ID       kernel.MessageID `json:"id"`
		SenderID kernel.UserID    `json:"sender_id"`
		RoomID   kernel.RoomID    `json:"room_id"`
		Content  string           `json:"content"`
		SentAt   time.Time        `json:"sent_at"`
	}

	return json.Marshal(&Alias{
		ID:       m.id,
		SenderID: m.senderID,
		RoomID:   m.roomID,
		Content:  m.content,
		SentAt:   m.sentAt,
	})
}

func (m *RoomMessage) GetEvents() []event.Event {
	return m.eventManager.GetEvents()
}
