package websocket

import (
	"encoding/json"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

var _ application.MessageNotifier = (*MessageNotifier)(nil)

const NotifyingMessageTopic websocket.ResponseTopic = "notification.notifying_message"

type ResponseData struct {
	ID      domain.MessageID
	Sender  kernel.UserID
	State   domain.MessageState
	Content string
	SentAt  time.Time
}

type MessageNotifier struct {
	manager *websocket.Manager
}

func NewMessageNotifier(manager *websocket.Manager) *MessageNotifier {
	return &MessageNotifier{manager: manager}
}

func (n *MessageNotifier) Notify(recipient kernel.UserID, message *domain.Message) error {
	responseData := &ResponseData{
		ID:      message.ID(),
		Sender:  message.Sender(),
		State:   message.State(),
		Content: message.Content(),
		SentAt:  message.SentAt(),
	}

	data, err := json.Marshal(responseData)
	if err != nil {
		return err
	}

	response := &websocket.Response{
		ResponseTopic: NotifyingMessageTopic,
		Data:          data,
	}

	payload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return n.manager.SendTo(recipient, payload)
}
