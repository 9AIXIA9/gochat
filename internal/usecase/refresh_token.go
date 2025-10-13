package usecase

import (
	"context"
	"gochat/internal/domain"
	"time"
)

type RefreshToken struct {
	expireDuration time.Duration
	domain.RefreshTokenAggregate
	domain.AuthTokenGenerator
	domain.RefreshTokenGenerator
}

func NewRefreshToken(expireDuration time.Duration,
	refreshTokenAggregate domain.RefreshTokenAggregate,
	authTokenGenerator domain.AuthTokenGenerator,
	refreshTokenGenerator domain.RefreshTokenGenerator) domain.RefreshTokenUsecase {
	return &RefreshToken{
		expireDuration:        expireDuration,
		RefreshTokenAggregate: refreshTokenAggregate,
		AuthTokenGenerator:    authTokenGenerator,
		RefreshTokenGenerator: refreshTokenGenerator,
	}
}

func (uc *RefreshToken) Execute(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.Response, error) {
	//查询 refresh token是否存在
	info, err := uc.FindRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return domain.UnauthorizedResponse, nil
	}

	//生成 refresh token
	rToken, err := uc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	//存储 refresh token
	if err = uc.SaveRefreshToken(ctx, rToken, info, uc.expireDuration); err != nil {
		return nil, err
	}

	//生成 auth token
	aToken, err := uc.GenerateAuthToken(info.Auth)
	if err != nil {
		return nil, err
	}
	return domain.NewSuccessResponse(domain.RefreshTokenResponse{
		AuthToken:    aToken,
		RefreshToken: rToken,
	}), nil
}
