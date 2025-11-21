package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type PrivateMessageRepository interface {
	PrivateMessageSaver
	PrivateMessageFinder
}

type PrivateMessageSaver interface {
	SavePrivateMessage(ctx context.Context, message *PrivateMessage) error
}

type PrivateMessageFinder interface {
	FindPrivateMessage(ctx context.Context, messageID kernel.MessageID) (*PrivateMessage, error)
}
