//go:generate mockgen -source=user_repository.go -destination=./mocks/mock_user_repository.go -package=mocks
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
