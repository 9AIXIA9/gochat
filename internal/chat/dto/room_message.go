package dto

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type RoomMessage struct {
	ID       kernel.MessageID                      `json:"id" example:"019b593b-462e-74d6-bfda-0e103a172190"`
	SenderID kernel.UserID                         `json:"sender_id" example:"019b593b-462e-74d6-bfda-0e103a172190"`
	States   map[kernel.UserID]domain.MessageState `json:"states"`
	RoomID   kernel.RoomID                         `json:"room_id" example:"019b593b-462e-74d6-bfda-0e103a172190"`
	Content  string                                `json:"content" example:"Hello-Gochat!"`
	SentAt   time.Time                             `json:"sent_at" example:"2025-12-26 05:38:19.740"`
}

func ToRoomMessageDTO(message *domain.RoomMessage) *RoomMessage {
	return &RoomMessage{
		ID:       message.ID(),
		SenderID: message.SenderID(),
		States:   message.States(),
		RoomID:   message.RoomID(),
		Content:  message.Content(),
		SentAt:   message.SentAt(),
	}
}

func ToRoomMessageDTOs(messages []*domain.RoomMessage) []*RoomMessage {
	dtos := make([]*RoomMessage, 0, len(messages))
	for _, message := range messages {
		dtos = append(dtos, ToRoomMessageDTO(message))
	}
	return dtos
}
