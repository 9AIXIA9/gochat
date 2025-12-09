package dto

import (
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type SystemMessage struct {
	ID          kernel.MessageID    `json:"id"`
	RecipientID kernel.UserID       `json:"recipient_id"`
	State       domain.MessageState `json:"state"`
	Content     string              `json:"content"`
	SentAt      time.Time           `json:"sent_at"`
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
