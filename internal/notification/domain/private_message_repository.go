package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

//TODO 保存和创建要分开

type PrivateMessageRepository interface {
	PrivateMessageSaver
	PrivateMessagesSaver
	PrivateMessageFinder
	UndeliveredPrivateMessageFinder
}

type PrivateMessageSaver interface {
	SavePrivateMessage(ctx context.Context, message *PrivateMessage) error
}

type PrivateMessagesSaver interface {
	SavePrivateMessages(ctx context.Context, messages []*PrivateMessage) error
}

type PrivateMessageFinder interface {
	FindPrivateMessage(ctx context.Context, messageID kernel.MessageID) (*PrivateMessage, error)
}

type UndeliveredPrivateMessageFinder interface {
	FindUndeliveredPrivateMessages(ctx context.Context, userID kernel.UserID) ([]*PrivateMessage, error)
}
