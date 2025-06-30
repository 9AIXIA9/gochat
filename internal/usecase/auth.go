package usecase

import (
	"context"
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

func (uc *Auth) ParseToken(ctx context.Context, tokenStr string) (*domain.AuthInfo, error) {
	return jwt.ParseToken(ctx, uc.conf.Secret, tokenStr)
}
