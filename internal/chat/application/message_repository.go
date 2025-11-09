package application

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
)

type MessageRepository interface {
	PrivateMessageSaver
	RoomMessageSaver
	MessageStateUpdater
}

type PrivateMessageSaver interface {
	SavePrivateMessage(ctx context.Context, recipient kernel.UserID, message *domain.Message) error
}

type RoomMessageSaver interface {
	SaveRoomMessage(ctx context.Context, members []kernel.UserID, messages *domain.Message) error
}

type MessageStateUpdater interface {
	UpdateState(ctx context.Context, messageID domain.MessageID, recipientID kernel.UserID, newState domain.MessageState) error
}
