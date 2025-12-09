package dto

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type PrivateMessage struct {
	ID          kernel.MessageID    `json:"id"`
	SenderID    kernel.UserID       `json:"sender_id"`
	RecipientID kernel.UserID       `json:"recipient_id"`
	State       domain.MessageState `json:"state"`
	Content     string              `json:"content"`
	SentAt      time.Time           `json:"sent_at"`
}

func ToPrivateMessageDTO(message *domain.PrivateMessage) *PrivateMessage {
	return &PrivateMessage{
		ID:          message.ID(),
		SenderID:    message.SenderID(),
		RecipientID: message.RecipientID(),
		State:       message.State(),
		Content:     message.Content(),
		SentAt:      message.SentAt(),
	}
}

func ToPrivateMessageDTOs(messages []*domain.PrivateMessage) []*PrivateMessage {
	dtos := make([]*PrivateMessage, 0, len(messages))
	for _, message := range messages {
		dtos = append(dtos, ToPrivateMessageDTO(message))
	}
	return dtos
}
