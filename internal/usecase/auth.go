package usecase

import (
	"gochat/internal/domain"
)

type Auth struct {
	domain.AuthTokenParser
}

func NewAuth(parser domain.AuthTokenParser) domain.AuthUsecase {
	return &Auth{
		AuthTokenParser: parser,
	}
}
