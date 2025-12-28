package dto

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type PrivateMessage struct {
	ID          kernel.MessageID    `json:"id" example:"019b593b-462e-74d6-bfda-0e103a172190"`
	SenderID    kernel.UserID       `json:"sender_id" example:"019b5929-65bc-7549-89c6-3f7dc6872577"`
	RecipientID kernel.UserID       `json:"recipient_id" example:"019b5929-5f6a-73aa-8adf-57cebe980725"`
	State       domain.MessageState `json:"state" example:"read"`
	Content     string              `json:"content" example:"hello gochat!"`
	SentAt      time.Time           `json:"sent_at" example:"2025-12-26 05:56:55.470"`
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
