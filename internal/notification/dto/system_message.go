package dto

import (
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type SystemMessage struct {
	ID          kernel.MessageID    `json:"id" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	RecipientID kernel.UserID       `json:"recipient_id" example:"019b593b-462e-74d6-bfda-0e103a172192"`
	State       domain.MessageState `json:"state" example:"delivered"`
	Content     string              `json:"content" example:"Hello-Gochat!"`
	SentAt      time.Time           `json:"sent_at" example:"2025-12-26 05:38:19.740"`
}

func ToSystemMessageDTO(message *domain.SystemMessage) *SystemMessage {
	return &SystemMessage{
		ID:          message.ID(),
		RecipientID: message.RecipientID(),
		State:       message.State(),
		Content:     message.Content(),
		SentAt:      message.SentAt(),
	}
}

func ToSystemMessageDTOs(messages []*domain.SystemMessage) []*SystemMessage {
	dtos := make([]*SystemMessage, 0, len(messages))
	for _, message := range messages {
		dtos = append(dtos, ToSystemMessageDTO(message))
	}
	return dtos
}
