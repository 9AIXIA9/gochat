package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type MessageRepository interface {
	PrivateMessageSaver
	PrivateMessageFinder
	RoomMessageSaver
	RoomMessageFinder
}

type PrivateMessageSaver interface {
	SavePrivateMessage(ctx context.Context, message *PrivateMessage) error
}

type PrivateMessageFinder interface {
	FindPrivateMessage(ctx context.Context, messageID kernel.MessageID) (*PrivateMessage, error)
}

type RoomMessageSaver interface {
	SaveRoomMessage(ctx context.Context, message *RoomMessage) error
}

type RoomMessageFinder interface {
	FindRoomMessage(ctx context.Context, messageID kernel.MessageID) (*RoomMessage, error)
}
