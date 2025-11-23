package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type UserRepository interface {
	UserCreator
	UserFinderByNumber
	UserFinderByID
}

type UserCreator interface {
	Create(ctx context.Context, user *User) error
}

type UserFinderByNumber interface {
	FindByNumber(ctx context.Context, number kernel.UserNumber) (*User, error)
}

type UserFinderByID interface {
	FindByID(ctx context.Context, id kernel.UserID) (*User, error)
}
