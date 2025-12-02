package application

import (
	"context"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type FriendshipCreatedUseCase kernel.UseCase[*FriendshipCreatedInput, *kernel.NoOutput]

type FriendshipCreatedInput struct {
	RequestID kernel.OperationID
	UserID1   kernel.UserID
	UserID2   kernel.UserID
}

func (r *FriendshipCreatedInput) Validate() error {
	if len(r.RequestID) == 0 || len(r.UserID1) == 0 || len(r.UserID2) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type friendshipCreatedUseCase struct {
	idGenerator event.IDGenerator
}

func NewFriendshipCreatedUseCase(
	idGenerator event.IDGenerator,
) (FriendshipCreatedUseCase, error) {
	if err := utils.CheckInterfaces(
		idGenerator,
	); err != nil {
		return nil, err
	}
	return &friendshipCreatedUseCase{
		idGenerator: idGenerator,
	}, nil
}

func (uc *friendshipCreatedUseCase) Execute(context.Context, *FriendshipCreatedInput) (*kernel.NoOutput, error) {
	//TODO 发送给通知上下文 让他来通知双方
	return nil, nil
}
