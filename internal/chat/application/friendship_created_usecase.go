package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type FriendshipCreatedUseCase kernel.UseCase[*FriendshipCreatedInput, *kernel.NoOutput]

type FriendshipCreatedInput struct {
	ID      domain.FriendshipID
	UserID1 kernel.UserID
	UserID2 kernel.UserID
}

func (r *FriendshipCreatedInput) Validate() error {
	if len(r.ID) == 0 || len(r.UserID1) == 0 || len(r.UserID2) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type friendshipCreatedUseCase struct {
	friendshipCreator domain.FriendshipCreator
}

func NewFriendshipCreatedUseCase(
	friendshipCreator domain.FriendshipCreator,
) (FriendshipCreatedUseCase, error) {
	if err := utils.CheckInterfaces(friendshipCreator); err != nil {
		return nil, err
	}

	return &friendshipCreatedUseCase{
		friendshipCreator: friendshipCreator,
	}, nil
}

func (uc *friendshipCreatedUseCase) Execute(ctx context.Context, input *FriendshipCreatedInput) (*kernel.NoOutput, error) {
	friendship := domain.CreateFriendship(input.ID, input.UserID1, input.UserID2)

	return nil, uc.friendshipCreator.Create(ctx, friendship)
}
