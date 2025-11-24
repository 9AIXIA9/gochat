package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomMessageRepository interface {
	RoomMessageCreator
	RoomMessageUpdater
	RoomMessagesUpdater
	RoomMessageFinderByID
	UserRoomMessagesFinderByState
}

type RoomMessageCreator interface {
	Create(ctx context.Context, message *RoomMessage) error
}

type RoomMessageUpdater interface {
	Update(ctx context.Context, message *RoomMessage) error
}

type RoomMessagesUpdater interface {
	Updates(ctx context.Context, messages []*RoomMessage) error
}

type RoomMessageFinderByID interface {
	FindRoomMessage(ctx context.Context, messageID kernel.MessageID) (*RoomMessage, error)
}

type UserRoomMessagesFinderByState interface {
	FindsByState(ctx context.Context, userID kernel.UserID, state MessageState) ([]*RoomMessage, error)
}
