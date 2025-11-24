package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type PrivateMessageRepository interface {
	PrivateMessageCreator
	PrivateMessageUpdater
	PrivateMessagesUpdater
	PrivateMessageFinderByID
	UserPrivateMessagesFinderByState
}

type PrivateMessageCreator interface {
	Create(ctx context.Context, message *PrivateMessage) error
}

type PrivateMessageUpdater interface {
	Update(ctx context.Context, message *PrivateMessage) error
}

type PrivateMessagesUpdater interface {
	Updates(ctx context.Context, messages []*PrivateMessage) error
}

type PrivateMessageFinderByID interface {
	FindByID(ctx context.Context, messageID kernel.MessageID) (*PrivateMessage, error)
}

type UserPrivateMessagesFinderByState interface {
	FindsByState(ctx context.Context, userID kernel.UserID, state MessageState) ([]*PrivateMessage, error)
}
