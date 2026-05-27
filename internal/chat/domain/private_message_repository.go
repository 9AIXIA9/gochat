//go:generate mockgen -source=private_message_repository.go -destination=./mocks/mock_private_message_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type PrivateMessageRepository interface {
	PrivateMessageCreator
	PrivateMessageFinder
	PrivateMessagesFinderByUserIDs
}

type PrivateMessageCreator interface {
	Create(ctx context.Context, message *PrivateMessage) error
}

type PrivateMessageFinder interface {
	FindPrivateMessage(ctx context.Context, messageID kernel.MessageID) (*PrivateMessage, error)
}

type PrivateMessagesFinderByUserIDs interface {
	FindsByUserIDs(ctx context.Context, userID1, userID2 kernel.UserID, limit int, baseID kernel.MessageID) ([]*PrivateMessage, error)
}
