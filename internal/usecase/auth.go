package usecase

import (
	"context"
	"gochat/internal/config"
	"gochat/internal/domain"
	"gochat/internal/infra/jwt"
)

type Auth struct {
	conf *config.Token
}

func NewAuth(conf *config.Token) domain.AuthUsecase {
	return &Auth{conf: conf}
}

func (uc *Auth) ParseAuthToken(ctx context.Context, token domain.AuthToken) (*domain.AuthInfo, error) {
	return jwt.ParseAuthToken(ctx, uc.conf.Auth.Secret, token)
}
