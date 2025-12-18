package websocket

import (
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

var _ domain.SystemMessageNotifier = (*SystemMessageNotifier)(nil)

const NotifySystemMessageTopic websocket.Topic = "notification.notify_system_message"

type NotifySystemMessageBody struct {
	ID      kernel.MessageID    `json:"id"`
	State   domain.MessageState `json:"state"`
	Content string              `json:"content"`
	SentAt  time.Time           `json:"sent_at"`
}

type SystemMessageNotifier struct {
	manager *websocket.Manager
}

func NewSystemMessageNotifier(manager *websocket.Manager) *SystemMessageNotifier {
	return &SystemMessageNotifier{manager: manager}
}

func (n *SystemMessageNotifier) Notify(message *domain.SystemMessage) error {

	return n.manager.SendTo(message.RecipientID(), &websocket.Message{
		Topic: NotifySystemMessageTopic,
		Body: &NotifySystemMessageBody{
			ID:      message.ID(),
			State:   message.State(),
			Content: message.Content(),
			SentAt:  message.SentAt(),
		},
	})
}
