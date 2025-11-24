package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type UserRepository interface {
	UserFinderByNumber
	UserCreator
}

type UserCreator interface {
	Create(ctx context.Context, user *User) error
}

type UserFinderByNumber interface {
	FindByNumber(ctx context.Context, number kernel.UserNumber) (*User, error)
}
