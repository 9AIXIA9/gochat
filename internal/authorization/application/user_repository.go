package application

import (
	"context"
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type UserRepository interface {
	UserSaver
	UserFinderByNumber
	UserLoggedInAtUpdater
	UserFinderByID
}

type UserSaver interface {
	Save(ctx context.Context, user *domain.User) error
}

type UserFinderByNumber interface {
	FindByNumber(ctx context.Context, number kernel.UserNumber) (*domain.User, error)
}

type UserLoggedInAtUpdater interface {
	UpdateLoggedInAt(ctx context.Context, userID kernel.UserID, time time.Time) error
}

type UserFinderByID interface {
	FindByID(ctx context.Context, id kernel.UserID) (*domain.User, error)
}
