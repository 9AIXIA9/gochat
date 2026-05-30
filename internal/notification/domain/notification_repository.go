//go:generate mockgen -source=notification_repository.go -destination=./mocks/mock_notification_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type NotificationRepository interface {
	NotificationCreator
	NotificationDeleter
}

type NotificationCreator interface {
	Create(ctx context.Context, notification *Notification) error
}

type NotificationDeleter interface {
	Delete(ctx context.Context, id kernel.MessageID) error
}
