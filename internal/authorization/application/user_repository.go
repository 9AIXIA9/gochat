package application

import (
	"context"
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type UserSaver interface {
	Save(ctx context.Context, user *domain.User) error
}

type UserFinderByNumber interface {
	FindByNumber(ctx context.Context, number domain.UserNumber) (*domain.User, error)
}

type UserLoggedInAtUpdater interface {
	UpdateLoggedInAt(ctx context.Context, userID kernel.UserID, time time.Time) error
}
