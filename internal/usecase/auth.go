package usecase

import (
	"gochat/internal/config"
	"gochat/internal/domain"
	"gochat/internal/infra/jwt"
)

type Auth struct {
	conf *config.JWT
}

func NewAuth(conf *config.JWT) domain.AuthUsecase {
	return &Auth{conf: conf}
}

func (uc *Auth) ParseToken(tokenStr string) (*domain.AuthInfo, error) {
	return jwt.ParseToken(uc.conf.Secret, tokenStr)
}
