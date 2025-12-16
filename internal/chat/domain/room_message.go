package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

type RoomMessage struct {
	id       kernel.MessageID
	senderID kernel.UserID
	states   map[kernel.UserID]MessageState
	roomID   kernel.RoomID
	content  string
	sentAt   time.Time

	eventManager *event.Manager
}

func LoadRoomMessage(
	id kernel.MessageID,
	senderID kernel.UserID,
	states map[kernel.UserID]MessageState,
	roomID kernel.RoomID,
	content string,
	sentAt time.Time,
) *RoomMessage {
	return &RoomMessage{
		id:           id,
		senderID:     senderID,
		states:       states,
		roomID:       roomID,
		content:      content,
		sentAt:       sentAt,
		eventManager: event.NewEventManager(),
	}
}

func CreateRoomMessage(
	roomID kernel.RoomID,
	senderID kernel.UserID,
	recipientIDs []kernel.UserID,
	content string,
	messageIDGenerator kernel.MessageIDGenerator,
	notifier RoomMessageNotifier,
) *RoomMessage {
	if len(recipientIDs) == 0 ||
		(len(recipientIDs) == 1 && recipientIDs[0] == senderID) ||
		recipientIDs == nil {
		//不用通知任何人
		return &RoomMessage{
			id:           messageIDGenerator.Generate(),
			senderID:     senderID,
			states:       nil,
			roomID:       roomID,
			content:      content,
			sentAt:       time.Now().UTC(),
			eventManager: event.NewEventManager(),
		}
	}

	states := make(map[kernel.UserID]MessageState, len(recipientIDs))
	for _, recipientID := range recipientIDs {
		if recipientID != senderID {
			states[recipientID] = MessageStateUndelivered
		}
	}

	message := &RoomMessage{
		id:           messageIDGenerator.Generate(),
		senderID:     senderID,
		states:       states,
		roomID:       roomID,
		content:      content,
		sentAt:       time.Now().UTC(),
		eventManager: event.NewEventManager(),
	}

	ids, err := notifier.Notify(message, recipientIDs)
	if err != nil {
		return message
	}

	for _, id := range ids {
		message.states[id] = MessageStateDelivered
	}

	return message
}

func (m *RoomMessage) Deliver(
	userID kernel.UserID,
	notifier RoomMessageNotifier,
) error {
	if state, ok := m.states[userID]; !ok || state != MessageStateUndelivered {
		return nil
	}

	ids, err := notifier.Notify(m, []kernel.UserID{userID})
	if err != nil {
		return err
	}

	for _, id := range ids {
		m.states[id] = MessageStateDelivered
	}
	return nil
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

func (m *RoomMessage) States() map[kernel.UserID]MessageState {
	return m.states
}

func (m *RoomMessage) State(id kernel.UserID) MessageState {
	return m.states[id]
}

func (m *RoomMessage) GetEvents() []event.Event {
	return m.eventManager.GetEvents()
}
