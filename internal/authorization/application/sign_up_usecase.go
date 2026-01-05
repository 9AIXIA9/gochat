package application

import (
	"context"
	"errors"
	"gochat/internal/authorization/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
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
}

func NewSignUpUseCase(
	eventIDGenerator event.IDGenerator,
	userIDGenerator domain.UserIDGenerator,
	numberGenerator domain.UserNumberGenerator,
	encryptor domain.Encryptor,
	userCreator domain.UserCreator,
) (SignUpUseCase, error) {
	if err := validate.NotNil(
		eventIDGenerator,
		userIDGenerator,
		numberGenerator,
		encryptor,
		userCreator,
	); err != nil {
		return nil, err
	}
	return &signUpUseCase{
		eventIDGenerator: eventIDGenerator,
		userIDGenerator:  userIDGenerator,
		numberGenerator:  numberGenerator,
		encryptor:        encryptor,
		userCreator:      userCreator,
	}, nil
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

	err = uc.userCreator.Create(ctx, user)
	if err != nil {
		if errors.Is(err, myErrors.ErrDuplicatedKey) {
			return nil, myErrors.WrapBusiness(err, "email is already used")
		}
		return nil, err
	}
	return &SignUpOutput{UserNumber: user.Number()}, nil
}
