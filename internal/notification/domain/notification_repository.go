//go:generate mockgen -source=notification_repository.go -destination=./mocks/mock_notification_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type NotificationRepository interface {
	NotificationCreator
	NotificationStateUpdater
}

type NotificationCreator interface {
	Create(ctx context.Context, notification *Notification) error
}

type NotificationStateUpdater interface {
	UpdateStateByID(ctx context.Context, id kernel.MessageID, state NotificationState) error
}
