package dto

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type RoomMessage struct {
	ID       kernel.MessageID                      `json:"id"`
	SenderID kernel.UserID                         `json:"sender_id"`
	States   map[kernel.UserID]domain.MessageState `json:"states"`
	RoomID   kernel.RoomID                         `json:"room_id"`
	Content  string                                `json:"content"`
	SentAt   time.Time                             `json:"sent_at"`
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
