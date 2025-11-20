package application

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"

	"context"
)

type RefreshAccessTokenUseCase kernel.UseCase[*RefreshAccessTokenInput, *RefreshAccessTokenOutput]

type RefreshAccessTokenInput struct {
	RefreshToken domain.RefreshToken
}

type RefreshAccessTokenOutput struct {
	AccessToken  domain.AccessToken
	RefreshToken *domain.RefreshTokenEntity
}

func (input *RefreshAccessTokenInput) Validate() error {
	return input.RefreshToken.Validate()
}

type refreshAccessTokenUseCase struct {
	refreshTokenSaver     domain.RefreshTokenSaver
	refreshTokenFinder    domain.RefreshTokenFinder
	accessTokenGenerator  domain.AccessTokenGenerator
	refreshTokenGenerator domain.RefreshTokenGenerator
}

func NewRefreshAccessTokenUseCase(
	refreshTokenSaver domain.RefreshTokenSaver,
	refreshTokenFinder domain.RefreshTokenFinder,
	accessTokenGenerator domain.AccessTokenGenerator,
	refreshTokenGenerator domain.RefreshTokenGenerator,
) RefreshAccessTokenUseCase {
	return &refreshAccessTokenUseCase{
		refreshTokenSaver:     refreshTokenSaver,
		refreshTokenFinder:    refreshTokenFinder,
		accessTokenGenerator:  accessTokenGenerator,
		refreshTokenGenerator: refreshTokenGenerator,
	}
}
func (uc *refreshAccessTokenUseCase) Execute(ctx context.Context, input *RefreshAccessTokenInput) (*RefreshAccessTokenOutput, error) {
	refreshToken, err := uc.refreshTokenFinder.FindByToken(ctx, input.RefreshToken)
	if err != nil {
		return nil, err
	}

	if err := refreshToken.Refresh(
		uc.refreshTokenGenerator,
	); err != nil {
		return nil, err
	}

	accessToken, err := refreshToken.GenerateAccessToken(
		uc.accessTokenGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.refreshTokenSaver.Save(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &RefreshAccessTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
