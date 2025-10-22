package useCase

import (
	"context"
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"
	"gochat/internal/shared/kernel/event"
)

var _ domain.LoginUseCase = (*login)(nil)

type login struct {
	idGenerator           kernel.IDGenerator
	hashComparator        kernel.HashComparator
	userFinder            domain.UserFinder
	refreshTokenSaver     domain.RefreshTokenSaver
	accessTokenGenerator  kernel.AccessTokenGenerator
	refreshTokenGenerator domain.RandomStringGenerator
	eventSaver            event.Saver
}

func NewLogin(
	idGenerator kernel.IDGenerator,
	hashComparator kernel.HashComparator,
	userFinder domain.UserFinder,
	refreshTokenSaver domain.RefreshTokenSaver,
	accessTokenGenerator kernel.AccessTokenGenerator,
	refreshTokenGenerator domain.RandomStringGenerator,
	eventSaver event.Saver,
) domain.LoginUseCase {
	return &login{
		idGenerator:           idGenerator,
		hashComparator:        hashComparator,
		userFinder:            userFinder,
		refreshTokenSaver:     refreshTokenSaver,
		accessTokenGenerator:  accessTokenGenerator,
		refreshTokenGenerator: refreshTokenGenerator,
		eventSaver:            eventSaver,
	}
}

func (uc *login) Execute(ctx context.Context, input *domain.LoginInput) (*domain.LoginOutput, error) {
	user, err := uc.userFinder.FindByNumber(ctx, input.Number)
	if err != nil {
		return nil, err
	}

	refreshToken, accessToken, err := user.Login(input.Password, uc.idGenerator, uc.hashComparator, uc.refreshTokenGenerator, uc.accessTokenGenerator)
	if err != nil {
		return nil, err
	}

	if err := uc.refreshTokenSaver.Save(ctx, refreshToken); err != nil {
		return nil, err
	}

	if err := uc.eventSaver.Save(ctx, user.GetEvents()); err != nil {
		return nil, err
	}

	return &domain.LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
