package useCase

import (
	"context"
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"
	"gochat/internal/shared/kernel/event"
)

type signUp struct {
	minPasswordLength int
	maxPasswordLength int
	idGenerator       kernel.IDGenerator
	numberGenerator   kernel.NumberGenerator
	hashEncryptor     kernel.HashEncryptor
	userSaver         domain.UserSaver
	eventSaver        event.Saver
}

func NewSignUp(
	idGenerator kernel.IDGenerator,
	numberGenerator kernel.NumberGenerator,
	hashEncryptor kernel.HashEncryptor,
	userSaver domain.UserSaver,
	eventSaver event.Saver,
) domain.SignUpUseCase {
	return &signUp{
		idGenerator:     idGenerator,
		numberGenerator: numberGenerator,
		hashEncryptor:   hashEncryptor,
		userSaver:       userSaver,
		eventSaver:      eventSaver,
	}
}

func (uc *signUp) Execute(ctx context.Context, input *domain.SignUpInput) (*domain.SignUpOutput, error) {
	user, err := domain.SignUp(input.Email, input.Password, uc.hashEncryptor, uc.idGenerator, uc.numberGenerator)
	if err != nil {
		return nil, err
	}

	err = uc.userSaver.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	err = uc.eventSaver.Save(ctx, user.GetEvents())
	if err != nil {
		return nil, err
	}

	return &domain.SignUpOutput{UserNumber: user.Number()}, nil
}
