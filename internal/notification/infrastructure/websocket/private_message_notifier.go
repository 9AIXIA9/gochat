package websocket

import (
	"encoding/json"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

var _ domain.PrivateMessageNotifier = (*PrivateMessageNotifier)(nil)

const NotifyPrivateMessageTopic websocket.Topic = "notification.notify_private_message"

type NotifyPrivateMessageResponseData struct {
	ID       kernel.MessageID    `json:"id"`
	SenderID kernel.UserID       `json:"sender_id"`
	State    domain.MessageState `json:"state"`
	Content  string              `json:"content"`
	SentAt   time.Time           `json:"sent_at"`
}

type PrivateMessageNotifier struct {
	manager *websocket.Manager
}

func NewPrivateMessageNotifier(manager *websocket.Manager) *PrivateMessageNotifier {
	return &PrivateMessageNotifier{manager: manager}
}

func (n *PrivateMessageNotifier) NotifyPrivateMessage(message *domain.PrivateMessage) error {
	responseData := &NotifyPrivateMessageResponseData{
		ID:       message.ID(),
		SenderID: message.SenderID(),
		State:    message.State(),
		Content:  message.Content(),
		SentAt:   message.SentAt(),
	}

	data, err := json.Marshal(responseData)
	if err != nil {
		return err
	}

	return n.manager.SendTo(message.RecipientID(), &websocket.Response{
		Topic: NotifyPrivateMessageTopic,
		Data:  data,
	})
}
