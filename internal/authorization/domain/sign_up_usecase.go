package domain

import (
	"gochat/internal/shared/kernel"
)

type SignUpUseCase kernel.UseCase[*SignUpInput, *SignUpOutput]

type SignUpInput struct {
	Email    kernel.Email
	Password Password
}

func (r *SignUpInput) Validate() error {
	return r.Password.Validate()
}

type SignUpOutput struct {
	UserNumber UserNumber
}
