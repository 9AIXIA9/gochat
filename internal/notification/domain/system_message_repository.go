//go:generate mockgen -source=system_message_repository.go -destination=./mocks/mock_system_message_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type SystemMessageRepository interface {
	SystemMessageCreator
	SystemMessagesStatesUpdaterByMessageIDs
	SystemMessageFinderByID
	UserSystemMessagesFinderByState
	SystemMessageFinderByUserID
}

type SystemMessagesStatesUpdaterByMessageIDs interface {
	UpdatesByMessageIDs(ctx context.Context, userID kernel.UserID, ids []kernel.MessageID, state MessageState) error
}

type SystemMessageCreator interface {
	Create(ctx context.Context, message *SystemMessage) error
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
