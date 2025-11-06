package application

import (
	"context"
	"gochat/internal/chat/domain"
)

type UserFinder interface {
	FindByNumber(ctx context.Context, number domain.UserNumber) (*domain.User, error)
}
