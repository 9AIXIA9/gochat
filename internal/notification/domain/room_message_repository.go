package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomMessageRepository interface {
	RoomMessageSaver
	RoomMessagesSaver
	RoomMessageFinder
	UndeliveredRoomMessageFinder
}

type RoomMessageSaver interface {
	SaveRoomMessage(ctx context.Context, message *RoomMessage) error
}

type RoomMessagesSaver interface {
	SaveRoomMessages(ctx context.Context, message []*RoomMessage) error
}

type RoomMessageFinder interface {
	FindRoomMessage(ctx context.Context, messageID kernel.MessageID) (*RoomMessage, error)
}

type UndeliveredRoomMessageFinder interface {
	FindUndeliveredRoomMessages(ctx context.Context, userID kernel.UserID) ([]*RoomMessage, error)
}
