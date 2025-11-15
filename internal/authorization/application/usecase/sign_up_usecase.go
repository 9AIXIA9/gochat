package usecase

import (
	"context"
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
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
	UserNumber domain.UserNumber
}

type signUpUseCase struct {
	eventIDGenerator event.IDGenerator
	userIDGenerator  application.UserIDGenerator
	numberGenerator  application.UserNumberGenerator
	encryptor        application.Encryptor
	userSaver        application.UserSaver
	eventSaver       event.UnpublishedSaver
}

func NewSignUpUseCase(
	eventIDGenerator event.IDGenerator,
	userIDGenerator application.UserIDGenerator,
	numberGenerator application.UserNumberGenerator,
	encryptor application.Encryptor,
	userSaver application.UserSaver,
	eventSaver event.UnpublishedSaver,
) SignUpUseCase {
	return &signUpUseCase{
		eventIDGenerator: eventIDGenerator,
		userIDGenerator:  userIDGenerator,
		numberGenerator:  numberGenerator,
		encryptor:        encryptor,
		userSaver:        userSaver,
		eventSaver:       eventSaver,
	}
}

func (uc *signUpUseCase) Execute(ctx context.Context, input *SignUpInput) (*SignUpOutput, error) {
	now := time.Now().UTC()
	passwordEncrypted, err := uc.encryptor.Encrypt(input.Password.String())
	if err != nil {
		return nil, err
	}

	user := domain.NewUser(uc.userIDGenerator.Generate(),
		input.Email,
		uc.numberGenerator.Generate(),
		passwordEncrypted,
		now,
		now,
	)

	err = user.SignUp(uc.eventIDGenerator)
	if err != nil {
		return nil, err
	}

	err = uc.userSaver.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	err = uc.eventSaver.Saves(ctx, user.GetEvents())
	if err != nil {
		return nil, err
	}

	return &SignUpOutput{UserNumber: user.Number()}, nil
}
