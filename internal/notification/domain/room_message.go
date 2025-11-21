package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

type RoomMessage struct {
	id           kernel.MessageID
	senderID     kernel.UserID
	roomID       kernel.RoomID
	recipientIDs []kernel.UserID
	states       map[kernel.UserID]MessageState
	content      string
	sentAt       time.Time
}

func LoadRoomMessage(
	id kernel.MessageID,
	senderID kernel.UserID,
	roomID kernel.RoomID,
	recipientIDs []kernel.UserID,
	states map[kernel.UserID]MessageState,
	content string,
	sentAt time.Time,
) *RoomMessage {
	return &RoomMessage{
		id:           id,
		senderID:     senderID,
		roomID:       roomID,
		recipientIDs: recipientIDs,
		states:       states,
		content:      content,
		sentAt:       sentAt,
	}
}

func ReceiveRoomMessage(
	id kernel.MessageID,
	senderID kernel.UserID,
	roomID kernel.RoomID,
	recipientIDs []kernel.UserID,
	content string,
	sentAt time.Time,
) *RoomMessage {
	states := make(map[kernel.UserID]MessageState)
	for _, recipientID := range recipientIDs {
		states[recipientID] = MessageStateReceived
	}

	return &RoomMessage{
		id:           id,
		senderID:     senderID,
		roomID:       roomID,
		recipientIDs: recipientIDs,
		states:       states,
		content:      content,
		sentAt:       sentAt,
	}
}

func (m *RoomMessage) Deliver(
	notifier RoomMessageNotifier,
) error {
	recipientIDsNotDelivered := make([]kernel.UserID, 0)
	for _, recipientID := range m.recipientIDs {
		if m.states[recipientID] == MessageStateReceived {
			recipientIDsNotDelivered = append(recipientIDsNotDelivered, recipientID)
		}
	}

	idsDelivered, err := notifier.NotifyRoomMessage(m, recipientIDsNotDelivered)
	if err != nil {
		return err
	}

	for _, idDelivered := range idsDelivered {
		m.states[idDelivered] = MessageStateDelivered
	}
	return nil
}

func (m *RoomMessage) Read(recipientID kernel.UserID) {
	m.states[recipientID] = MessageStateDelivered
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

func (m *RoomMessage) States() map[kernel.UserID]MessageState {
	return m.states
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
