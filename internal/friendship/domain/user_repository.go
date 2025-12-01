//go:generate mockgen -source=user_repository.go -destination=./mocks/mock_user_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type UserRepository interface {
	UserCreator
	UserSaver
	UserFinderByID
	UserFinderByNumber
}

type UserCreator interface {
	Create(ctx context.Context, user *User) error
}

type UserSaver interface {
	Save(ctx context.Context, user *User) error
}

type UserFinderByID interface {
	FindByID(ctx context.Context, id kernel.UserID) (*User, error)
}

type UserFinderByNumber interface {
	FindByNumber(ctx context.Context, number kernel.UserNumber) (*User, error)
}
