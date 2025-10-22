package converters

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/authorization/infrastructure/persistence/models"
	"gochat/internal/shared/kernel"
)

var _ kernel.GenericModelConverter[*models.User, *domain.User] = (*UserConverter)(nil)

type UserConverter struct {
}

func (c *UserConverter) ToModel(user *domain.User) (*models.User, error) {
	return &models.User{
		ID:             user.ID(),
		Email:          user.Email(),
		Number:         user.Number(),
		PasswordHash:   user.PasswordHash(),
		SignedUpAt:     user.SignedUpAt(),
		LastLoggedInAt: user.LastLoggedInAt(),
	}, nil
}

func (c *UserConverter) ToDomain(user *models.User) (*domain.User, error) {
	return domain.NewUser(user.ID, user.Email, user.Number, user.PasswordHash, user.SignedUpAt, user.LastLoggedInAt), nil
}
