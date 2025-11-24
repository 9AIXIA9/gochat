package domain

import (
	"context"
)

type UserRepository interface {
	UserCreator
}

type UserCreator interface {
	Create(ctx context.Context, user *User) error
}
