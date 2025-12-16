package application

import (
	"context"
	"gochat/internal/friendship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

const (
	defaultFriendshipsLimit = 50
	maxFriendshipsLimit     = 100
)

type ListFriendshipsUseCase kernel.UseCase[*ListFriendshipsInput, *ListFriendshipsOutput]

type ListFriendshipsInput struct {
	UserID kernel.UserID
	BaseID domain.FriendshipID
	Limit  int
}

func (r *ListFriendshipsInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id can't be empty")
	}

	return nil
}

type ListFriendshipsOutput struct {
	Friendships []*domain.Friendship
}

type listFriendshipsUseCase struct {
	friendshipFinder domain.FriendshipsFinderByUserID
}

func NewListFriendshipsUseCase(
	friendshipFinder domain.FriendshipsFinderByUserID,
) (ListFriendshipsUseCase, error) {
	if err := utils.CheckInterfaces(friendshipFinder); err != nil {
		return nil, err
	}
	return &listFriendshipsUseCase{
		friendshipFinder: friendshipFinder,
	}, nil
}

func (uc *listFriendshipsUseCase) Execute(ctx context.Context, input *ListFriendshipsInput) (*ListFriendshipsOutput, error) {
	limit := input.Limit
	switch {
	case limit <= 0:
		limit = defaultFriendshipsLimit
	case limit > maxFriendshipsLimit:
		limit = maxFriendshipsLimit
	}

	friendships, err := uc.friendshipFinder.FindsByUserID(ctx, input.UserID, limit, input.BaseID)
	if err != nil {
		return nil, err
	}

	return &ListFriendshipsOutput{
		Friendships: friendships,
	}, nil
}
