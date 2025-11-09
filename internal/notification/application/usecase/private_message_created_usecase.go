package usecase

import (
	"context"
	"gochat/internal/shared/kernel"
)

type PrivateMessageCreatedUseCase kernel.UseCase[*PrivateMessageCreatedInput, *kernel.NoOutput]

type PrivateMessageCreatedInput struct {
}

func (r *PrivateMessageCreatedInput) Validate() error {
	return nil
}

type privateMessageCreatedUseCase struct {
}

func NewPrivateMessageCreatedUseCase() PrivateMessageCreatedUseCase {
	return &privateMessageCreatedUseCase{}
}

func (uc *privateMessageCreatedUseCase) Execute(context.Context, *PrivateMessageCreatedInput) (*kernel.NoOutput, error) {
	return nil, nil
}
