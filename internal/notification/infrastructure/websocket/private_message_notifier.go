package websocket

import (
	"encoding/json"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

//TODO 断线重连后消息通知未发送问题

var _ domain.PrivateMessageNotifier = (*PrivateMessageNotifier)(nil)

const NotifyPrivateMessageTopic websocket.ResponseTopic = "notification.notify_private_message"

type NotifyPrivateMessageResponseData struct {
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

	response := &websocket.Response{
		ResponseTopic: NotifyPrivateMessageTopic,
		Data:          data,
	}

	payload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return n.manager.SendTo(message.RecipientID(), payload)
}
