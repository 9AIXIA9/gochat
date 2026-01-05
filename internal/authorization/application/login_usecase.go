package application

import (
	"context"
	"errors"
	"gochat/internal/authorization/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type LoginUseCase kernel.UseCase[*LoginInput, *LoginOutput]

type LoginInput struct {
	Number   kernel.UserNumber
	Password domain.Password
}

type LoginOutput struct {
	AccessToken  domain.AccessToken
	RefreshToken *domain.RefreshTokenEntity
}

func (i *LoginInput) Validate() error {
	if err := i.Number.Validate(); err != nil {
		return err
	}
	return i.Password.Validate()
}

type loginUseCase struct {
	comparator            domain.Comparator
	userFinder            domain.UserFinderByNumber
	refreshTokenUpserter  domain.RefreshTokenUpserter
	accessTokenGenerator  domain.AccessTokenGenerator
	refreshTokenGenerator domain.RefreshTokenGenerator
}

func NewLoginUseCase(
	comparator domain.Comparator,
	userFinder domain.UserFinderByNumber,
	refreshTokenUpserter domain.RefreshTokenUpserter,
	accessTokenGenerator domain.AccessTokenGenerator,
	refreshTokenGenerator domain.RefreshTokenGenerator,
) (LoginUseCase, error) {
	if err := validate.NotNil(
		comparator,
		userFinder,
		refreshTokenUpserter,
		accessTokenGenerator,
		refreshTokenGenerator,
	); err != nil {
		return nil, err
	}
	return &loginUseCase{
		comparator:            comparator,
		userFinder:            userFinder,
		refreshTokenUpserter:  refreshTokenUpserter,
		accessTokenGenerator:  accessTokenGenerator,
		refreshTokenGenerator: refreshTokenGenerator,
	}, nil
}

func (uc *loginUseCase) Execute(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	user, err := uc.userFinder.FindByNumber(ctx, input.Number)
	if err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return nil, domain.ErrInvalidPassword
		}
		return nil, err
	}

	if err := user.PasswordEncrypted().Compare(input.Password, uc.comparator); err != nil {
		return nil, err
	}

	refreshToken, err := domain.CreateRefreshToken(
		user.ID(),
		uc.refreshTokenGenerator,
	)
	if err != nil {
		return nil, err
	}

	accessToken, err := refreshToken.GenerateAccessToken(uc.accessTokenGenerator)
	if err != nil {
		return nil, err
	}

	if err := uc.refreshTokenUpserter.Upsert(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
