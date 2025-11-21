package websocket

import (
	"encoding/json"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

var _ domain.RoomMessageNotifier = (*RoomMessageNotifier)(nil)

const NotifyingRoomMessageTopic websocket.ResponseTopic = "notification.notifying_room_message"

type NotifyingRoomMessageResponseData struct {
	ID       kernel.MessageID
	SenderID kernel.UserID
	RoomID   kernel.RoomID
	States   map[kernel.UserID]domain.MessageState
	Content  string
	SentAt   time.Time
}

type RoomMessageNotifier struct {
	manager *websocket.Manager
}

func NewRoomMessageNotifier(manager *websocket.Manager) *RoomMessageNotifier {
	return &RoomMessageNotifier{manager: manager}
}

func (n *RoomMessageNotifier) NotifyRoomMessage(message *domain.RoomMessage, recipients []kernel.UserID) ([]kernel.UserID, error) {
	responseData := &NotifyingRoomMessageResponseData{
		ID:       message.ID(),
		SenderID: message.SenderID(),
		RoomID:   message.RoomID(),
		States:   message.States(),
		Content:  message.Content(),
		SentAt:   message.SentAt(),
	}

	data, err := json.Marshal(responseData)
	if err != nil {
		return nil, err
	}

	response := &websocket.Response{
		ResponseTopic: NotifyingRoomMessageTopic,
		Data:          data,
	}

	payload, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return n.manager.Broadcast(recipients, payload), nil
}
