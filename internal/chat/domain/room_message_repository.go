//go:generate mockgen -source=room_message_repository.go -destination=./mocks/mock_room_message_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomMessageRepository interface {
	RoomMessageCreator
	RoomMessagesUpdater
	RoomMessageFinder
	RoomMessagesFinderByRecipientIDAndState
	RoomMessagesStatesUpdaterByUserIDAndRoomID
	RoomMessagesStatesUpdaterByMessageIDs
	RoomMessagesFinderByRoomIDAndUserID
}

type RoomMessageCreator interface {
	Create(ctx context.Context, message *RoomMessage) error
}

type RoomMessagesUpdater interface {
	Updates(ctx context.Context, messages []*RoomMessage) error
}

type RoomMessageFinder interface {
	FindRoomMessage(ctx context.Context, messageID kernel.MessageID) (*RoomMessage, error)
}

type RoomMessagesFinderByRecipientIDAndState interface {
	FindRoomMessagesByRecipientIDAndState(ctx context.Context, recipientID kernel.UserID, state MessageState, limit int) ([]*RoomMessage, error)
}

type RoomMessagesStatesUpdaterByUserIDAndRoomID interface {
	UpdatesByUserIDAndRoomID(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID, state MessageState) error
}

type RoomMessagesStatesUpdaterByMessageIDs interface {
	UpdatesByMessageIDs(ctx context.Context, userID kernel.UserID, ids []kernel.MessageID, state MessageState) error
}

type RoomMessagesFinderByRoomIDAndUserID interface {
	FindsByRoomIDAndUserID(ctx context.Context, roomID kernel.RoomID, userID kernel.UserID, limit int, baseID kernel.MessageID) ([]*RoomMessage, error)
}
