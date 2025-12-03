package websocket

import (
	"encoding/json"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

var _ domain.RoomMessageNotifier = (*RoomMessageNotifier)(nil)

const NotifyRoomMessageTopic websocket.Topic = "notification.notify_room_message"

type NotifyRoomMessageResponseData struct {
	ID       kernel.MessageID                      `json:"id"`
	SenderID kernel.UserID                         `json:"sender_id"`
	RoomID   kernel.RoomID                         `json:"room_id"`
	States   map[kernel.UserID]domain.MessageState `json:"states"`
	Content  string                                `json:"content"`
	SentAt   time.Time                             `json:"sent_at"`
}

type RoomMessageNotifier struct {
	manager *websocket.Manager
}

func NewRoomMessageNotifier(manager *websocket.Manager) *RoomMessageNotifier {
	return &RoomMessageNotifier{manager: manager}
}

func (n *RoomMessageNotifier) Notify(message *domain.RoomMessage, recipients []kernel.UserID) ([]kernel.UserID, error) {
	responseData := &NotifyRoomMessageResponseData{
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

	return n.manager.Broadcast(recipients, &websocket.Response{
		Topic: NotifyRoomMessageTopic,
		Data:  data,
	},
	), nil
}
