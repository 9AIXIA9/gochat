//go:generate mockgen -source=system_message_repository.go -destination=./mocks/mock_system_message_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type SystemMessageRepository interface {
	SystemMessageCreator
	SystemMessageUpdater
	SystemMessagesUpdater
	SystemMessageFinderByID
	UserSystemMessagesFinderByState
	SystemMessageFinderByUserID
}

type SystemMessageCreator interface {
	Create(ctx context.Context, message *SystemMessage) error
}

type SystemMessageUpdater interface {
	Update(ctx context.Context, message *SystemMessage) error
}

type SystemMessagesUpdater interface {
	Updates(ctx context.Context, messages []*SystemMessage) error
}

type SystemMessageFinderByID interface {
	FindByID(ctx context.Context, messageID kernel.MessageID) (*SystemMessage, error)
}

type UserSystemMessagesFinderByState interface {
	FindsByState(ctx context.Context, userID kernel.UserID, state MessageState, limit int) ([]*SystemMessage, error)
}

type SystemMessageFinderByUserID interface {
	FindsByUserID(ctx context.Context, userID kernel.UserID, limit int, baseID kernel.MessageID) ([]*SystemMessage, error)
}
