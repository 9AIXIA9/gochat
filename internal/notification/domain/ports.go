//go:generate mockgen -source=ports.go -destination=./mocks/mock_ports.go -package=mocks
package domain

import (
	"gochat/internal/shared/kernel"
)

type PrivateMessageNotifier interface {
	NotifyPrivateMessage(message *PrivateMessage) error
}

type RoomMessageNotifier interface {
	NotifyRoomMessage(message *RoomMessage, recipients []kernel.UserID) ([]kernel.UserID, error)
}
