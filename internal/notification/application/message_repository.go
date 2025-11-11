package application

import (
	"context"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
)

type MessageRepository interface {
	MessageSaver
	MessageFinder
	MessageStatesUpdater
}

type MessageSaver interface {
	SaveMessage(ctx context.Context, recipient kernel.UserID, message *domain.Message) error
}

type MessageFinder interface {
	FindMessagesByUserID(ctx context.Context, userID kernel.UserID) ([]*domain.Message, error)
}

type MessageStatesUpdater interface {
	UpdateMessageStates(ctx context.Context, userID kernel.UserID, messageIDs []domain.MessageID, newState domain.MessageState) error
}
