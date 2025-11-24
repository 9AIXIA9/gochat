package application

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"
)

type AccessTokenParser interface {
	Parse(token domain.AccessToken) (kernel.UserID, error)
}
