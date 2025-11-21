package websocket

import (
	"encoding/json"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

var _ domain.PrivateMessageNotifier = (*PrivateMessageNotifier)(nil)

const NotifyingPrivateMessageTopic websocket.ResponseTopic = "notification.notifying_private_message"

type NotifyingPrivateMessageResponseData struct {
	ID       kernel.MessageID
	SenderID kernel.UserID
	State    domain.MessageState
	Content  string
	SentAt   time.Time
}

type PrivateMessageNotifier struct {
	manager *websocket.Manager
}

func NewPrivateMessageNotifier(manager *websocket.Manager) *PrivateMessageNotifier {
	return &PrivateMessageNotifier{manager: manager}
}

func (n *PrivateMessageNotifier) NotifyPrivateMessage(message *domain.PrivateMessage) error {
	responseData := &NotifyingPrivateMessageResponseData{
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

	response := &websocket.Response{
		ResponseTopic: NotifyingPrivateMessageTopic,
		Data:          data,
	}

	payload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return n.manager.SendTo(message.RecipientID(), payload)
}
