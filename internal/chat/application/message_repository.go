package application

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
)

type PrivateMessagesSaver interface {
	SavePrivateMessages(ctx context.Context, messages []*domain.PrivateMessage) error
}

type MessageStateUpdater interface {
	Update(ctx context.Context, messageID domain.MessageID, recipientID kernel.UserID, newState domain.MessageState) error
}
