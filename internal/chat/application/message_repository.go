package application

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
)

type MessageRepository interface {
	PrivateMessageSaver
	PrivateMessageFinder
	RoomMessageSaver
	RoomMessageFinder
}

type PrivateMessageSaver interface {
	SavePrivateMessage(ctx context.Context, recipient kernel.UserID, message *domain.Message) error
}

type PrivateMessageFinder interface {
	FindPrivateMessage(ctx context.Context, recipient kernel.UserID, messageID kernel.MessageID) (*domain.Message, error)
}

type RoomMessageSaver interface {
	SaveRoomMessage(ctx context.Context, roomID kernel.RoomID, messages *domain.Message) error
}

type RoomMessageFinder interface {
	FindRoomMessage(ctx context.Context, roomID kernel.RoomID, messageID kernel.MessageID) (*domain.Message, error)
}
