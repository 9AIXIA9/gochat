package application

import (
	"errors"
	"gochat/internal/authorization/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"

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
	refreshTokenUpserter  domain.RefreshTokenUpserter
	refreshTokenFinder    domain.RefreshTokenFinder
	accessTokenGenerator  domain.AccessTokenGenerator
	refreshTokenGenerator domain.RefreshTokenGenerator
}

func NewRefreshAccessTokenUseCase(
	refreshTokenUpserter domain.RefreshTokenUpserter,
	refreshTokenFinder domain.RefreshTokenFinder,
	accessTokenGenerator domain.AccessTokenGenerator,
	refreshTokenGenerator domain.RefreshTokenGenerator,
) (RefreshAccessTokenUseCase, error) {
	if err := validate.NotNil(
		refreshTokenUpserter,
		refreshTokenFinder,
		accessTokenGenerator,
		refreshTokenGenerator,
	); err != nil {
		return nil, err
	}
	return &refreshAccessTokenUseCase{
		refreshTokenUpserter:  refreshTokenUpserter,
		refreshTokenFinder:    refreshTokenFinder,
		accessTokenGenerator:  accessTokenGenerator,
		refreshTokenGenerator: refreshTokenGenerator,
	}, nil
}

func (uc *refreshAccessTokenUseCase) Execute(ctx context.Context, input *RefreshAccessTokenInput) (*RefreshAccessTokenOutput, error) {
	refreshToken, err := uc.refreshTokenFinder.FindByToken(ctx, input.RefreshToken)
	if err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return nil, domain.ErrInvalidRefreshToken
		}
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

	if err := uc.refreshTokenUpserter.Upsert(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &RefreshAccessTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
