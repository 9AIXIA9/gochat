package application

import (
	"context"
	"errors"
	"gochat/internal/friendship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type AgreeFriendRequestUseCase kernel.UseCase[*AgreeFriendRequestInput, *kernel.NoOutput]

type AgreeFriendRequestInput struct {
	UserID    kernel.UserID
	RequestID kernel.OperationID
}

func (r *AgreeFriendRequestInput) Validate() error {
	if len(r.RequestID) == 0 || len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "request id and user id can't be empty")
	}

	return nil
}

type agreeFriendRequestUseCase struct {
	friendRequestFinderByRequestID domain.FriendRequestFinderByID
	friendRequestUpdater           domain.FriendRequestUpdater
	idGenerator                    event.IDGenerator
}

func NewAgreeFriendRequestUseCase(
	friendRequestFinderByRequestID domain.FriendRequestFinderByID,
	friendRequestUpdater domain.FriendRequestUpdater,
	idGenerator event.IDGenerator,
) (AgreeFriendRequestUseCase, error) {
	if err := validate.NotNil(
		friendRequestFinderByRequestID, friendRequestUpdater, idGenerator,
	); err != nil {
		return nil, err
	}
	return &agreeFriendRequestUseCase{
		friendRequestFinderByRequestID: friendRequestFinderByRequestID,
		friendRequestUpdater:           friendRequestUpdater,
		idGenerator:                    idGenerator,
	}, nil
}

func (uc *agreeFriendRequestUseCase) Execute(ctx context.Context, input *AgreeFriendRequestInput) (*kernel.NoOutput, error) {
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

	if err := req.Agree(uc.idGenerator); err != nil {
		return nil, err
	}

	if err := uc.friendRequestUpdater.Update(ctx, req); err != nil {
		return nil, err
	}

	return nil, nil
}
