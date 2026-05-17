package domain

import (
	"context"
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
	ctx context.Context,
	roomID kernel.RoomID,
	senderID kernel.UserID,
	content string,
	finder RoomshipsFinderByRoomID,
	messageIDGenerator kernel.MessageIDGenerator,
	notifier RoomMessageNotifier,
) (*RoomMessage, error) {
	roomships, err := finder.FindsByRoomID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	message, err := createRoomMessageByRoomships(roomships, roomID, senderID, content, messageIDGenerator, notifier)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func createRoomMessageByRoomships(
	roomships []*Roomship,
	roomID kernel.RoomID,
	senderID kernel.UserID,
	content string,
	messageIDGenerator kernel.MessageIDGenerator,
	notifier RoomMessageNotifier,
) (*RoomMessage, error) {
	if len(roomships) == 0 {
		return nil, ErrRoomNotFound
	}

	if len(content) == 0 {
		return nil, ErrEmptyMessageContent
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
			states:       nil,
			roomID:       roomID,
			content:      content,
			sentAt:       time.Now().UTC(),
			eventManager: event.NewEventManager(),
		}, nil
	}

	states := make(map[kernel.UserID]MessageState, len(recipientIDs))
	for _, recipientID := range recipientIDs {
		states[recipientID] = MessageStateUndelivered
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

	_, err := notifier.Notify(message, recipientIDs)
	if err != nil {
		return message, nil
	}

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

func (m *RoomMessage) States() map[kernel.UserID]MessageState {
	return m.states
}

func (m *RoomMessage) State(id kernel.UserID) MessageState {
	return m.states[id]
}

func (m *RoomMessage) UndeliveredRecipientIDs() []kernel.UserID {
	recipientIDs := make([]kernel.UserID, 0)
	for recipientID, state := range m.states {
		if state == MessageStateUndelivered {
			recipientIDs = append(recipientIDs, recipientID)
		}
	}
	return recipientIDs
}

func (m *RoomMessage) GetEvents() []event.Event {
	return m.eventManager.GetEvents()
}
