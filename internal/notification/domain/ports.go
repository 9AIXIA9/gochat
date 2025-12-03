//go:generate mockgen -source=ports.go -destination=./mocks/mock_ports.go -package=mocks
package domain

import (
	"gochat/internal/shared/kernel"
)

type SystemMessageNotifier interface {
	Notify(message *SystemMessage) error
}

type PrivateMessageNotifier interface {
	Notify(message *PrivateMessage) error
}

type RoomMessageNotifier interface {
	Notify(message *RoomMessage, recipients []kernel.UserID) ([]kernel.UserID, error)
}
