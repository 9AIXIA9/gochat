package application

import (
	"context"
	"gochat/internal/friendship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type ListFriendRequestsUseCase kernel.UseCase[*ListFriendRequestsInput, *ListFriendRequestsOutput]

const (
	defaultFriendRequestsLimit = 10
	maxFriendRequestsLimit     = 100
)

type ListFriendRequestsInput struct {
	UserID kernel.UserID
	BaseID kernel.OperationID
	Limit  int
}

func (r *ListFriendRequestsInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id can't be empty")
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
	limit := input.Limit
	switch {
	case limit <= 0:
		limit = defaultFriendRequestsLimit
	case limit > maxFriendRequestsLimit:
		limit = maxFriendRequestsLimit
	}

	requests, err := uc.finder.FindsByUserID(ctx, input.UserID, limit, input.BaseID)
	if err != nil {
		return nil, err
	}

	return &ListFriendRequestsOutput{
		Requests: requests,
	}, nil
}
