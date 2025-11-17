package application

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
)

type UserRepository interface {
	UserFinder
	UserNumberSaver
}

type UserFinder interface {
	FindByNumber(ctx context.Context, number kernel.UserNumber) (*domain.User, error)
}

type UserNumberSaver interface {
	SaveNumber(ctx context.Context, userID kernel.UserID, number kernel.UserNumber) error
}
