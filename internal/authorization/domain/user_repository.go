package domain

import (
	"context"
)

type UserRepository interface {
	UserSaver
	UserFinder
}

type UserSaver interface {
	Save(ctx context.Context, user *User) error
}

type UserFinder interface {
	FindByNumber(ctx context.Context, number UserNumber) (*User, error)
}
