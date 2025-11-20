package uuid

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"

	"github.com/google/uuid"
)

var _ domain.UserIDGenerator = (*UserIDGenerator)(nil)

type UserIDGenerator struct{}

func NewUserIDGenerator() *UserIDGenerator {
	return &UserIDGenerator{}
}

func (UserIDGenerator) Generate() kernel.UserID {
	return kernel.UserID(uuid.Must(uuid.NewV7()).String())
}
