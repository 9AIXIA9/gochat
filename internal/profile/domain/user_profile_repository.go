//go:generate mockgen -source=user_profile_repository.go -destination=./mocks/mock_user_profile_repository.go -package=mocks
package domain

import "context"

type UserProfileRepository interface {
	UserProfileCreator
}

type UserProfileCreator interface {
	Create(ctx context.Context, profile *UserProfile) error
}
