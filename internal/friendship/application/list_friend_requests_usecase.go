package application

import (
	"context"
	"gochat/internal/friendship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type ListFriendRequestsUseCase kernel.UseCase[*ListFriendRequestsInput, *ListFriendRequestsOutput]

const defaultLimit = 10

type ListFriendRequestsInput struct {
	UserID kernel.UserID
	BaseID kernel.OperationID
	Limit  int
}

func (r *ListFriendRequestsInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type ListFriendRequestsOutput struct {
	Requests []*domain.FriendRequest
}

type listFriendRequestsUseCase struct {
	finder domain.FriendRequestsFinderByUserID
}

func NewListFriendRequestsUseCase(
	finder domain.FriendRequestsFinderByUserID,
) (ListFriendRequestsUseCase, error) {
	if err := utils.CheckInterfaces(finder); err != nil {
		return nil, err
	}
	return &listFriendRequestsUseCase{
		finder: finder,
	}, nil
}

func (uc *listFriendRequestsUseCase) Execute(ctx context.Context, input *ListFriendRequestsInput) (*ListFriendRequestsOutput, error) {
	if input.Limit == 0 {
		input.Limit = defaultLimit
	}

	requests, err := uc.finder.FindsByUserID(ctx, input.UserID, input.BaseID, input.Limit)
	if err != nil {
		return nil, err
	}

	return &ListFriendRequestsOutput{
		Requests: requests,
	}, nil
}
