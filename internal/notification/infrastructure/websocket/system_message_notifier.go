package websocket

import (
	"encoding/json"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

var _ domain.SystemMessageNotifier = (*SystemMessageNotifier)(nil)

const NotifySystemMessageTopic websocket.Topic = "notification.notify_system_message"

type NotifySystemMessageResponseData struct {
	ID      kernel.MessageID          `json:"id"`
	Topic   domain.SystemMessageTopic `json:"topic"`
	State   domain.MessageState       `json:"state"`
	Content []byte                    `json:"content"`
	SentAt  time.Time                 `json:"sent_at"`
}

type SystemMessageNotifier struct {
	manager *websocket.Manager
}

func NewSystemMessageNotifier(manager *websocket.Manager) *SystemMessageNotifier {
	return &SystemMessageNotifier{manager: manager}
}

func (n *SystemMessageNotifier) Notify(message *domain.SystemMessage) error {
	responseData := &NotifySystemMessageResponseData{
		ID:      message.ID(),
		Topic:   message.Topic(),
		State:   message.State(),
		Content: message.Content(),
		SentAt:  message.SentAt(),
	}

	data, err := json.Marshal(responseData)
	if err != nil {
		return err
	}

	return n.manager.SendTo(message.RecipientID(), &websocket.Response{
		Topic: NotifySystemMessageTopic,
		Data:  data,
	})
}
