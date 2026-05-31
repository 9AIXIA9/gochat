package application

import (
	"context"
	"errors"
	"gochat/internal/authorization/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type LoginByEmailUseCase kernel.UseCase[*LoginByEmailInput, *LoginByEmailOutput]

type LoginByEmailInput struct {
	Email    kernel.Email
	Password domain.Password
}

type LoginByEmailOutput struct {
	AccessToken  domain.AccessToken
	RefreshToken *domain.RefreshTokenEntity
}

func (i *LoginByEmailInput) Validate() error {
	if err := i.Email.Validate(); err != nil {
		return err
	}
	return i.Password.Validate()
}

type loginByEmailUseCase struct {
	comparator            domain.Comparator
	userFinder            domain.UserFinderByEmail
	refreshTokenUpserter  domain.RefreshTokenUpserter
	accessTokenGenerator  domain.AccessTokenGenerator
	refreshTokenGenerator domain.RefreshTokenGenerator
}

func NewLoginByEmailUseCase(
	comparator domain.Comparator,
	userFinder domain.UserFinderByEmail,
	refreshTokenUpserter domain.RefreshTokenUpserter,
	accessTokenGenerator domain.AccessTokenGenerator,
	refreshTokenGenerator domain.RefreshTokenGenerator,
) (LoginByEmailUseCase, error) {
	if err := validate.NotNil(
		comparator,
		userFinder,
		refreshTokenUpserter,
		accessTokenGenerator,
		refreshTokenGenerator,
	); err != nil {
		return nil, err
	}
	return &loginByEmailUseCase{
		comparator:            comparator,
		userFinder:            userFinder,
		refreshTokenUpserter:  refreshTokenUpserter,
		accessTokenGenerator:  accessTokenGenerator,
		refreshTokenGenerator: refreshTokenGenerator,
	}, nil
}

func (uc *loginByEmailUseCase) Execute(ctx context.Context, input *LoginByEmailInput) (*LoginByEmailOutput, error) {
	user, err := uc.userFinder.FindByEmail(ctx, input.Email)
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

	return &LoginByEmailOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
