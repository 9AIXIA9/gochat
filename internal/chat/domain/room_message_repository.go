//go:generate mockgen -source=room_message_repository.go -destination=./mocks/mock_room_message_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomMessageRepository interface {
	RoomMessageCreator
	RoomMessageFinder
	RoomMessagesFinderByRoomIDAndUserIDWithoutRecipients
}

type RoomMessageCreator interface {
	Create(ctx context.Context, message *RoomMessage) error
}

type RoomMessageFinder interface {
	FindRoomMessage(ctx context.Context, messageID kernel.MessageID) (*RoomMessage, error)
}

type RoomMessagesFinderByRoomIDAndUserIDWithoutRecipients interface {
	FindsByRoomIDAndUserIDWithoutRecipients(ctx context.Context, roomID kernel.RoomID, userID kernel.UserID, limit int, baseID kernel.MessageID) ([]*RoomMessage, error)
}
