//go:generate mockgen -source=user_profile_repository.go -destination=./mocks/mock_user_profile_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type UserProfileRepository interface {
	UserProfileCreator
	UserProfileUpdater
	UserProfileFinder
}

type UserProfileCreator interface {
	Create(ctx context.Context, profile *UserProfile) error
}

type UserProfileUpdater interface {
	Update(ctx context.Context, profile *UserProfile) error
}

type UserProfileFinder interface {
	FindByID(ctx context.Context, id kernel.UserID) (*UserProfile, error)
}
