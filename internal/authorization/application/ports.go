package application

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"
)

type Encryptor interface {
	Encrypt(origin string) (encrypted string, err error)
}

type Comparator interface {
	Compare(encrypted string, origin string) error
}

type AccessTokenParser interface {
	Parse(token domain.AccessToken) (kernel.UserID, error)
}

type AccessTokenGenerator interface {
	Generate(userID kernel.UserID) (domain.AccessToken, error)
}

type RefreshTokenGenerator interface {
	Generate() (domain.RefreshToken, error)
}

type UserNumberGenerator interface {
	Generate() domain.UserNumber
}

type UserIDGenerator interface {
	Generate() kernel.UserID
}
