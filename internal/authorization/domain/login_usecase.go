package domain

import (
	"gochat/internal/shared/kernel"
)

type LoginUseCase kernel.UseCase[*LoginInput, *LoginOutput]

type LoginInput struct {
	Number   UserNumber
	Password Password
}

type LoginOutput struct {
	AccessToken  kernel.AccessToken
	RefreshToken *RefreshToken
}

func (r *LoginInput) Validate() error {
	if err := r.Number.Validate(); err != nil {
		return err
	}
	return nil
}
