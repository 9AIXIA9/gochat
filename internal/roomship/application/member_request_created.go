package application

import (
	"context"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type MemberRequestCreatedUseCase kernel.UseCase[*MemberRequestCreatedInput, *kernel.NoOutput]

type MemberRequestCreatedInput struct {
	UserID kernel.UserID
}

func (r *MemberRequestCreatedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type memberRequestCreatedUseCase struct {
}

func NewMemberRequestCreatedUseCase() MemberRequestCreatedUseCase {
	return &memberRequestCreatedUseCase{}
}

func (uc *memberRequestCreatedUseCase) Execute(ctx context.Context, input *MemberRequestCreatedInput) (*kernel.NoOutput, error) {
	//TODO 通知上下文进行通知 各个admin
	return nil, nil
}
