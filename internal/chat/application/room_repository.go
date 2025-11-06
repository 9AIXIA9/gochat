package application

import (
	"context"
	"gochat/internal/chat/domain"
)

type RoomFinder interface {
	FindByNumber(ctx context.Context, number domain.RoomNumber) (*domain.Room, error)
}
