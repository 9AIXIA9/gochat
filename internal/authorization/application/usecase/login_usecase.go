package usecase

import (
	"context"
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type LoginUseCase kernel.UseCase[*LoginInput, *LoginOutput]

type LoginInput struct {
	Number   domain.UserNumber
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
	comparator            application.Comparator
	userFinder            application.UserFinderByNumber
	userUpdater           application.UserLoggedInAtUpdater
	refreshTokenSaver     application.RefreshTokenSaver
	accessTokenGenerator  application.AccessTokenGenerator
	refreshTokenGenerator application.RefreshTokenGenerator
}

func NewLoginUseCase(
	idGenerator event.IDGenerator,
	comparator application.Comparator,
	userFinder application.UserFinderByNumber,
	userUpdater application.UserLoggedInAtUpdater,
	refreshTokenSaver application.RefreshTokenSaver,
	accessTokenGenerator application.AccessTokenGenerator,
	refreshTokenGenerator application.RefreshTokenGenerator,
) LoginUseCase {
	return &loginUseCase{
		eventIDGenerator:      idGenerator,
		comparator:            comparator,
		userFinder:            userFinder,
		userUpdater:           userUpdater,
		refreshTokenSaver:     refreshTokenSaver,
		accessTokenGenerator:  accessTokenGenerator,
		refreshTokenGenerator: refreshTokenGenerator,
	}
}

func (uc *loginUseCase) Execute(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	user, err := uc.userFinder.FindByNumber(ctx, input.Number)
	if err != nil {
		return nil, err
	}

	if err := uc.comparator.Compare(user.PasswordEncrypted(), input.Password.String()); err != nil {
		return nil, err
	}

	user.Login()

	accessToken, err := uc.accessTokenGenerator.Generate(user.ID())
	if err != nil {
		return nil, err
	}

	token, err := uc.refreshTokenGenerator.Generate()
	if err != nil {
		return nil, err
	}

	refreshToken := domain.CreateRefreshToken(token, user.ID())

	if err := uc.refreshTokenSaver.Save(ctx, refreshToken); err != nil {
		return nil, err
	}

	if err := uc.userUpdater.UpdateLoggedInAt(ctx, user.ID(), user.LastLoggedInAt()); err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
