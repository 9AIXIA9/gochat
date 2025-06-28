package usecase

import (
	"context"
	"gochat/internal/config"
	"gochat/internal/domain"

	"github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	conf *config.JWT
}

func NewAuth(conf *config.JWT) domain.AuthUsecase {
	return &Auth{conf: conf}
}

func (uc *Auth) ParseToken(ctx context.Context, tokenStr string) bool {
	//解析 Token
	token, err := jwt.ParseWithClaims(tokenStr, &domain.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(uc.conf.Secret), nil
	})

	if err != nil {
		return false
	}

	//验证合理性
	claims, ok := token.Claims.(*domain.JwtCustomClaims)

	// 将用户信息存储到上下文中
	context.WithValue(ctx, domain.AuthKey, claims.Auth)
	return ok && token.Valid
}
