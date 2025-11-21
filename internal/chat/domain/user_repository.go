package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type UserRepository interface {
	UserFinderByNumber
	UserNumberSaver
}

type UserFinderByNumber interface {
	FindByNumber(ctx context.Context, number kernel.UserNumber) (*User, error)
}

type UserNumberSaver interface {
	SaveNumber(ctx context.Context, userID kernel.UserID, number kernel.UserNumber) error
}
