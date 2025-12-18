package websocket

import (
	"gochat/internal/chat/domain"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/kernel"
	"time"
)

var _ domain.PrivateMessageNotifier = (*PrivateMessageNotifier)(nil)

const NotifyPrivateMessageTopic websocket.Topic = "chat.notify_private_message"

type NotifyPrivateMessageBody struct {
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

func (n *PrivateMessageNotifier) Notify(message *domain.PrivateMessage) error {
	return n.manager.SendTo(message.RecipientID(), &websocket.Message{
		Topic: NotifyPrivateMessageTopic,
		Body: &NotifyPrivateMessageBody{
			ID:       message.ID(),
			SenderID: message.SenderID(),
			State:    message.State(),
			Content:  message.Content(),
			SentAt:   message.SentAt(),
		},
	})
}
