//go:generate mockgen -source=private_message_repository.go -destination=./mocks/mock_private_message_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type PrivateMessageRepository interface {
	PrivateMessageCreator
	PrivateMessagesUpdater
	PrivateMessageFinder
	PrivateMessagesFinderByRecipientIDAndState
}

type PrivateMessageCreator interface {
	Create(ctx context.Context, message *PrivateMessage) error
}

type PrivateMessagesUpdater interface {
	Updates(ctx context.Context, messages []*PrivateMessage) error
}

type PrivateMessageFinder interface {
	FindPrivateMessage(ctx context.Context, messageID kernel.MessageID) (*PrivateMessage, error)
}

type PrivateMessagesFinderByRecipientIDAndState interface {
	FindPrivateMessagesByRecipientIDAndState(ctx context.Context, recipientID kernel.UserID, state MessageState, limit int) ([]*PrivateMessage, error)
}
