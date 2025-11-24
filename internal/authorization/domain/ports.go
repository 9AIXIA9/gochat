package domain

import (
	"gochat/internal/shared/kernel"
)

type Encryptor interface {
	Encrypt(origin string) (encrypted string, err error)
}

type Comparator interface {
	Compare(encrypted string, origin string) error
}

type AccessTokenGenerator interface {
	Generate(userID kernel.UserID) (AccessToken, error)
}

type RefreshTokenGenerator interface {
	Generate() (RefreshToken, error)
}

type UserNumberGenerator interface {
	Generate() kernel.UserNumber
}

type UserIDGenerator interface {
	Generate() kernel.UserID
}
