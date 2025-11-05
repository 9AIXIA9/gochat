package application

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
)

type MessageSaver interface {
	Saves(ctx context.Context, messages []*domain.Message) error
}

type MessageStateUpdater interface {
	Update(ctx context.Context, messageID domain.MessageID, recipientID kernel.UserID, newState domain.MessageState) error
}
