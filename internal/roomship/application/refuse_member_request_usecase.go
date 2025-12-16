package application

import (
	"context"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type RefuseMemberRequestUseCase kernel.UseCase[*RefuseMemberRequestInput, *kernel.NoOutput]

type RefuseMemberRequestInput struct {
	UserID    kernel.UserID
	RequestID kernel.OperationID
}

func (r *RefuseMemberRequestInput) Validate() error {
	if len(r.UserID) == 0 || len(r.RequestID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id or request id is empty")
	}
	return nil
}

type refuseMemberRequestUseCase struct {
	requestFinder  domain.MemberRequestFinderByID
	roomshipFinder domain.RoomshipFinderByUserIDAndRoomID
	updater        domain.MemberRequestUpdater
}

func NewRefuseMemberRequestUseCase(
	requestFinder domain.MemberRequestFinderByID,
	roomshipFinder domain.RoomshipFinderByUserIDAndRoomID,
	updater domain.MemberRequestUpdater,
) (RefuseMemberRequestUseCase, error) {
	if err := utils.CheckInterfaces(
		requestFinder,
		roomshipFinder,
		updater,
	); err != nil {
		return nil, err
	}

	return &refuseMemberRequestUseCase{
		requestFinder:  requestFinder,
		roomshipFinder: roomshipFinder,
		updater:        updater,
	}, nil
}

func (uc *refuseMemberRequestUseCase) Execute(ctx context.Context, input *RefuseMemberRequestInput) (*kernel.NoOutput, error) {
	req, err := uc.requestFinder.FindByID(ctx, input.RequestID)
	if err != nil {
		return nil, err
	}

	roomship, err := uc.roomshipFinder.FindByUserIDAndRoomID(ctx, input.UserID, req.RoomID())
	if err != nil {
		return nil, err
	}

	if roomship == nil || !roomship.IsAdmin() {
		return nil, domain.ErrNotAdmin
	}

	if err := req.Refuse(input.UserID); err != nil {
		return nil, err
	}

	if err := uc.updater.Update(ctx, req); err != nil {
		return nil, err
	}

	return nil, nil
}
