//go:generate mockgen -source=user_profile_repository.go -destination=./mocks/mock_user_profile_repository.go -package=mocks
package domain

type UserProfileRepository interface {
}
