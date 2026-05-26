//go:generate mockgen -source=ports.go -destination=./mocks/mock_ports.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type PrivateMessageNotifier interface {
	Notify(ctx context.Context, message *PrivateMessage) error
}

type RoomMessageNotifier interface {
	Notify(ctx context.Context, message *RoomMessage, recipients []kernel.UserID) error
}
