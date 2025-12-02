package application

import (
	"context"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type FriendRequestCreatedUseCase kernel.UseCase[*FriendRequestCreatedInput, *kernel.NoOutput]

type FriendRequestCreatedInput struct {
	RequestID kernel.OperationID
}

func (r *FriendRequestCreatedInput) Validate() error {
	if len(r.RequestID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type friendRequestCreatedUseCase struct {
	idGenerator event.IDGenerator
}

func NewFriendRequestCreatedUseCase(
	idGenerator event.IDGenerator,
) (FriendRequestCreatedUseCase, error) {
	if err := utils.CheckInterfaces(
		idGenerator,
	); err != nil {
		return nil, err
	}
	return &friendRequestCreatedUseCase{
		idGenerator: idGenerator,
	}, nil
}

func (uc *friendRequestCreatedUseCase) Execute(context.Context, *FriendRequestCreatedInput) (*kernel.NoOutput, error) {
	//TODO 发送给通知上下文 让他来通知双方
	return nil, nil
}
