package application

import (
	"context"
	"gochat/internal/authorization/domain"
)

type UserSaver interface {
	Save(ctx context.Context, user *domain.User) error
}

type UserFinder interface {
	FindByNumber(ctx context.Context, number domain.UserNumber) (*domain.User, error)
}
