package application

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
)

type MessageRepository interface {
	PrivateMessageSaver
	RoomMessageSaver
}

type PrivateMessageSaver interface {
	SavePrivateMessage(ctx context.Context, recipient kernel.UserID, message *domain.Message) error
}

type RoomMessageSaver interface {
	SaveRoomMessage(ctx context.Context, roomID domain.RoomID, messages *domain.Message) error
}
