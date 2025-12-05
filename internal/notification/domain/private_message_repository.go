//go:generate mockgen -source=private_message_repository.go -destination=./mocks/mock_private_message_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type PrivateMessageRepository interface {
	PrivateMessageSaver
	PrivateMessageUpdater
	PrivateMessagesUpdater
	PrivateMessageFinderByID
	UserPrivateMessagesFinderByState
}

type PrivateMessageSaver interface {
	Save(ctx context.Context, message *PrivateMessage) error
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
