package application

import (
	"context"
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
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
	eventIDGenerator      event.IDGenerator
	comparator            domain.Comparator
	userFinder            domain.UserFinderByNumber
	refreshTokenCreator   domain.RefreshTokenCreator
	accessTokenGenerator  domain.AccessTokenGenerator
	refreshTokenGenerator domain.RefreshTokenGenerator
}

func NewLoginUseCase(
	idGenerator event.IDGenerator,
	comparator domain.Comparator,
	userFinder domain.UserFinderByNumber,
	refreshTokenCreator domain.RefreshTokenCreator,
	accessTokenGenerator domain.AccessTokenGenerator,
	refreshTokenGenerator domain.RefreshTokenGenerator,
) LoginUseCase {
	return &loginUseCase{
		eventIDGenerator:      idGenerator,
		comparator:            comparator,
		userFinder:            userFinder,
		refreshTokenCreator:   refreshTokenCreator,
		accessTokenGenerator:  accessTokenGenerator,
		refreshTokenGenerator: refreshTokenGenerator,
	}
}

func (uc *loginUseCase) Execute(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	user, err := uc.userFinder.FindByNumber(ctx, input.Number)
	if err != nil {
		return nil, err
	}

	refreshToken, err := user.Login(
		input.Password,
		uc.comparator,
		uc.refreshTokenGenerator,
	)
	if err != nil {
		return nil, err
	}

	accessToken, err := refreshToken.GenerateAccessToken(uc.accessTokenGenerator)
	if err != nil {
		return nil, err
	}

	if err := uc.refreshTokenCreator.Create(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
