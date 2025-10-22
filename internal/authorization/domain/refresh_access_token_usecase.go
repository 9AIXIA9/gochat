package domain

import (
	"gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type RefreshAccessTokenUseCase kernel.UseCase[*RefreshAccessTokenInput, *RefreshAccessTokenOutput]

type RefreshAccessTokenInput struct {
	RefreshTokenString string
}

type RefreshAccessTokenOutput struct {
	AccessToken  kernel.AccessToken
	RefreshToken *RefreshToken
}

func (input *RefreshAccessTokenInput) Validate() error {
	if len(input.RefreshTokenString) == 0 {
		return errors.ErrEmptyInput
	}
	return nil
}
