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
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "friendship id and user ids can't be empty")
	}

	return nil
}

type friendshipCreatedUseCase struct {
	friendshipSaver domain.FriendshipSaver
}

func NewFriendshipCreatedUseCase(
	friendshipSaver domain.FriendshipSaver,
) (FriendshipCreatedUseCase, error) {
	if err := utils.CheckInterfaces(friendshipSaver); err != nil {
		return nil, err
	}

	return &friendshipCreatedUseCase{
		friendshipSaver: friendshipSaver,
	}, nil
}

func (uc *friendshipCreatedUseCase) Execute(ctx context.Context, input *FriendshipCreatedInput) (*kernel.NoOutput, error) {
	friendship := domain.LoadFriendship(input.ID, input.UserID1, input.UserID2)

	return nil, uc.friendshipSaver.Save(ctx, friendship)
}
