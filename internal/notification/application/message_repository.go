package application

import (
	"context"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
)

type MessageRepository interface {
	MessageSaver
}

type MessageSaver interface {
	SaveMessage(ctx context.Context, recipient kernel.UserID, message *domain.Message) error
}
