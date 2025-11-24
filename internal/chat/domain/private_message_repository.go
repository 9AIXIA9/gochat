package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type PrivateMessageRepository interface {
	PrivateMessageCreator
	PrivateMessageFinder
}

type PrivateMessageCreator interface {
	Create(ctx context.Context, message *PrivateMessage) error
}

type PrivateMessageFinder interface {
	FindPrivateMessage(ctx context.Context, messageID kernel.MessageID) (*PrivateMessage, error)
}
