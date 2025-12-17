package application

import (
	"context"
	"errors"
	"gochat/internal/friendship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type RefuseFriendRequestUseCase kernel.UseCase[*RefuseFriendRequestInput, *kernel.NoOutput]

type RefuseFriendRequestInput struct {
	UserID    kernel.UserID
	RequestID kernel.OperationID
}

func (r *RefuseFriendRequestInput) Validate() error {
	if len(r.RequestID) == 0 || len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id and request id can't be empty")
	}

	return nil
}

type refuseFriendRequestUseCase struct {
	friendRequestFinderByRequestID domain.FriendRequestFinderByID
	friendRequestUpdater           domain.FriendRequestUpdater
}

func NewRefuseFriendRequestUseCase(
	friendRequestFinderByRequestID domain.FriendRequestFinderByID,
	friendRequestUpdater domain.FriendRequestUpdater,
) (RefuseFriendRequestUseCase, error) {
	if err := utils.CheckInterfaces(
		friendRequestFinderByRequestID, friendRequestUpdater,
	); err != nil {
		return nil, err
	}
	return &refuseFriendRequestUseCase{
		friendRequestFinderByRequestID: friendRequestFinderByRequestID,
		friendRequestUpdater:           friendRequestUpdater,
	}, nil
}

func (uc *refuseFriendRequestUseCase) Execute(ctx context.Context, input *RefuseFriendRequestInput) (*kernel.NoOutput, error) {
	req, err := uc.friendRequestFinderByRequestID.FindByID(ctx, input.RequestID)
	if err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return nil, myErrors.WrapBusiness(err, "request doesn't exist")
		}
		return nil, err
	}

	if req.To() != input.UserID {
		return nil, domain.ErrFriendRequestNotForUser
	}

	if err := req.Refuse(); err != nil {
		return nil, err
	}

	if err := uc.friendRequestUpdater.Update(ctx, req); err != nil {
		return nil, err
	}

	return nil, nil
}
