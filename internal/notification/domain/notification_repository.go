//go:generate mockgen -source=notification_repository.go -destination=./mocks/mock_notification_repository.go -package=mocks
package domain

import (
	"context"
)

type NotificationRepository interface {
	NotificationCreator
}

type NotificationCreator interface {
	Create(ctx context.Context, notification *Notification) error
}
