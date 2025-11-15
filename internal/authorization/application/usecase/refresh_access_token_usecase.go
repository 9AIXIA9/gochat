package usecase

import (
	"gochat/internal/authorization/application"
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
	refreshTokenSaver     application.RefreshTokenSaver
	refreshTokenFinder    application.RefreshTokenFinder
	accessTokenGenerator  application.AccessTokenGenerator
	refreshTokenGenerator application.RefreshTokenGenerator
}

func NewRefreshAccessTokenUseCase(
	refreshTokenSaver application.RefreshTokenSaver,
	refreshTokenFinder application.RefreshTokenFinder,
	accessTokenGenerator application.AccessTokenGenerator,
	refreshTokenGenerator application.RefreshTokenGenerator,
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

	if err := refreshToken.CanBeRefreshed(); err != nil {
		return nil, err
	}

	accessToken, err := uc.accessTokenGenerator.Generate(refreshToken.UserID())
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := uc.refreshTokenGenerator.Generate()
	if err != nil {
		return nil, err
	}

	refreshToken.RefreshAccessToken(newRefreshToken)

	if err := uc.refreshTokenSaver.Save(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &RefreshAccessTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
