package application

import (
	"context"
	"errors"
	"gochat/internal/authorization/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type LoginByNumberUseCase kernel.UseCase[*LoginByNumberInput, *LoginByNumberOutput]

type LoginByNumberInput struct {
	Number   kernel.UserNumber
	Password domain.Password
}

type LoginByNumberOutput struct {
	AccessToken  domain.AccessToken
	RefreshToken *domain.RefreshTokenEntity
}

func (i *LoginByNumberInput) Validate() error {
	if err := i.Number.Validate(); err != nil {
		return err
	}
	return i.Password.Validate()
}

type loginByNumberUseCase struct {
	comparator            domain.Comparator
	userFinder            domain.UserFinderByNumber
	refreshTokenUpserter  domain.RefreshTokenUpserter
	accessTokenGenerator  domain.AccessTokenGenerator
	refreshTokenGenerator domain.RefreshTokenGenerator
}

func NewLoginByNumberUseCase(
	comparator domain.Comparator,
	userFinder domain.UserFinderByNumber,
	refreshTokenUpserter domain.RefreshTokenUpserter,
	accessTokenGenerator domain.AccessTokenGenerator,
	refreshTokenGenerator domain.RefreshTokenGenerator,
) (LoginByNumberUseCase, error) {
	if err := validate.NotNil(
		comparator,
		userFinder,
		refreshTokenUpserter,
		accessTokenGenerator,
		refreshTokenGenerator,
	); err != nil {
		return nil, err
	}
	return &loginByNumberUseCase{
		comparator:            comparator,
		userFinder:            userFinder,
		refreshTokenUpserter:  refreshTokenUpserter,
		accessTokenGenerator:  accessTokenGenerator,
		refreshTokenGenerator: refreshTokenGenerator,
	}, nil
}

func (uc *loginByNumberUseCase) Execute(ctx context.Context, input *LoginByNumberInput) (*LoginByNumberOutput, error) {
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

	return &LoginByNumberOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
