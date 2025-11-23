package application

import (
	"context"
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type SignUpUseCase kernel.UseCase[*SignUpInput, *SignUpOutput]

type SignUpInput struct {
	Email    kernel.Email
	Password domain.Password
}

func (r *SignUpInput) Validate() error {
	if err := r.Password.Validate(); err != nil {
		return err
	}
	return r.Email.Validate()
}

type SignUpOutput struct {
	UserNumber kernel.UserNumber
}

type signUpUseCase struct {
	eventIDGenerator event.IDGenerator
	userIDGenerator  domain.UserIDGenerator
	numberGenerator  domain.UserNumberGenerator
	encryptor        domain.Encryptor
	userCreator      domain.UserCreator
	eventsCreator    event.UnpublishedEventsCreator
	unitOfWork       kernel.UnitOfWork
}

func NewSignUpUseCase(
	eventIDGenerator event.IDGenerator,
	userIDGenerator domain.UserIDGenerator,
	numberGenerator domain.UserNumberGenerator,
	encryptor domain.Encryptor,
	userCreator domain.UserCreator,
	eventsCreator event.UnpublishedEventsCreator,
	unitOfWork kernel.UnitOfWork,
) SignUpUseCase {
	return &signUpUseCase{
		eventIDGenerator: eventIDGenerator,
		userIDGenerator:  userIDGenerator,
		numberGenerator:  numberGenerator,
		encryptor:        encryptor,
		userCreator:      userCreator,
		eventsCreator:    eventsCreator,
		unitOfWork:       unitOfWork,
	}
}

func (uc *signUpUseCase) Execute(ctx context.Context, input *SignUpInput) (*SignUpOutput, error) {
	passwordEncrypted, err := input.Password.Encrypt(uc.encryptor)
	if err != nil {
		return nil, err
	}

	user, err := domain.CreateUser(
		input.Email,
		passwordEncrypted,
		uc.userIDGenerator,
		uc.numberGenerator,
		uc.eventIDGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.unitOfWork.Execute(ctx, func(txCtx context.Context) error {
		err = uc.userCreator.Create(txCtx, user)
		if err != nil {
			return err
		}

		err = uc.eventsCreator.CreateUnpublishedEvents(txCtx, user.GetEvents())
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &SignUpOutput{UserNumber: user.Number()}, nil
}
