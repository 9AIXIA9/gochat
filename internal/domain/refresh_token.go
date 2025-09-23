package domain

import (
	"context"
)

type RefreshToken string

type RefreshInfo struct {
	Auth *AuthInfo
}

type RefreshTokenUsecase interface {
	Execute(ctx context.Context, req *RefreshTokenRequest) (*Response, error)
	QueryRefreshToken(ctx context.Context, token RefreshToken) (*RefreshInfo, error)
	SaveRefreshToken(ctx context.Context, token RefreshToken, info *RefreshInfo) error
	GenerateRefreshToken() (RefreshToken, error)
	GenerateAuthToken(ctx context.Context, authInfo *AuthInfo) (AuthToken, error)
}
