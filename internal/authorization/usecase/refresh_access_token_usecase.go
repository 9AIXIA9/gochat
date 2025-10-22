package useCase

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"

	"context"
)

var _ domain.RefreshAccessTokenUseCase = (*refreshAccessToken)(nil)

type refreshAccessToken struct {
	refreshTokenSaver     domain.RefreshTokenSaver
	refreshTokenFinder    domain.RefreshTokenFinder
	accessTokenGenerator  kernel.AccessTokenGenerator
	refreshTokenGenerator domain.RandomStringGenerator
}

func NewRefreshAccessToken(
	refreshTokenSaver domain.RefreshTokenSaver,
	refreshTokenFinder domain.RefreshTokenFinder,
	accessTokenGenerator kernel.AccessTokenGenerator,
	refreshTokenGenerator domain.RandomStringGenerator,
) domain.RefreshAccessTokenUseCase {
	return &refreshAccessToken{
		refreshTokenSaver:     refreshTokenSaver,
		refreshTokenFinder:    refreshTokenFinder,
		accessTokenGenerator:  accessTokenGenerator,
		refreshTokenGenerator: refreshTokenGenerator,
	}
}
func (uc *refreshAccessToken) Execute(ctx context.Context, input *domain.RefreshAccessTokenInput) (*domain.RefreshAccessTokenOutput, error) {
	refreshToken, err := uc.refreshTokenFinder.FindByString(ctx, input.RefreshTokenString)
	if err != nil {
		return nil, err
	}

	accessToken, err := refreshToken.RefreshAccessToken(uc.refreshTokenGenerator, uc.accessTokenGenerator)
	if err != nil {
		return nil, err
	}

	if err := uc.refreshTokenSaver.Save(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &domain.RefreshAccessTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
