package kernel

import (
	"gochat/internal/shared/errors"
)

type ParseAccessTokenUseCase UseCase[*ParseAccessTokenInput, *ParseAccessTokenOutput]

type ParseAccessTokenInput struct {
	AccessToken AccessToken
}

func (input *ParseAccessTokenInput) Validate() error {
	if len(input.AccessToken) == 0 {
		return errors.ErrEmptyInput
	}
	return nil
}

type ParseAccessTokenOutput struct {
	UserID UserID
}
