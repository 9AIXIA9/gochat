package application

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
)

type PrivateMessageSaver interface {
	SavePrivateMessages(ctx context.Context, recipient kernel.UserID, message []*domain.Message) error
}

type RoomMessageSaver interface {
	SaveRoomMessages(ctx context.Context, roomID domain.RoomID, message []*domain.Message) error
}

type MessageStateUpdater interface {
	Update(ctx context.Context, messageID domain.MessageID, recipientID kernel.UserID, newState domain.MessageState) error
}
