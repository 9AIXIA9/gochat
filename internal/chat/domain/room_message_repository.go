package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomMessageRepository interface {
	RoomMessageSaver
	RoomMessageFinder
}

type RoomMessageSaver interface {
	SaveRoomMessage(ctx context.Context, message *RoomMessage) error
}

type RoomMessageFinder interface {
	FindRoomMessage(ctx context.Context, messageID kernel.MessageID) (*RoomMessage, error)
}
