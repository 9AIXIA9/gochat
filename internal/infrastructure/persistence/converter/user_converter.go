package converter

import (
	"gochat/internal/authorization/domain"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/infrastructure/persistence/model"
)

var _ gormutils.GenericModelConverter[*model.User, *domain.User] = (*UserConverter)(nil)

type UserConverter struct {
}

func (c *UserConverter) ToModel(user *domain.User) *model.User {
	return &model.User{
		ID:                user.ID(),
		Email:             user.Email(),
		Number:            user.Number(),
		PasswordEncrypted: user.PasswordEncrypted(),
		SignedUpAt:        user.SignedUpAt(),
		LastLoggedInAt:    user.LastLoggedInAt(),
	}
}

func (c *UserConverter) ToDomain(user *model.User) *domain.User {
	return domain.NewUser(user.ID, user.Email, user.Number, user.PasswordEncrypted, user.SignedUpAt, user.LastLoggedInAt)
}
