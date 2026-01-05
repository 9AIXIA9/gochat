package application

import (
	"context"
	"errors"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type AgreeMemberRequestUseCase kernel.UseCase[*AgreeMemberRequestInput, *kernel.NoOutput]

type AgreeMemberRequestInput struct {
	UserID    kernel.UserID
	RequestID kernel.OperationID
}

func (r *AgreeMemberRequestInput) Validate() error {
	if len(r.UserID) == 0 || len(r.RequestID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id or request id is empty")
	}
	return nil
}

type agreeMemberRequestUseCase struct {
	requestFinder    domain.MemberRequestFinderByID
	roomshipFinder   domain.RoomshipFinderByUserIDAndRoomID
	updater          domain.MemberRequestUpdater
	eventIDGenerator event.IDGenerator
}

func NewAgreeMemberRequestUseCase(
	requestFinder domain.MemberRequestFinderByID,
	roomshipFinder domain.RoomshipFinderByUserIDAndRoomID,
	updater domain.MemberRequestUpdater,
	eventIDGenerator event.IDGenerator,
) (AgreeMemberRequestUseCase, error) {
	if err := validate.NotNil(requestFinder,
		roomshipFinder,
		updater,
		eventIDGenerator,
	); err != nil {
		return nil, err
	}

	return &agreeMemberRequestUseCase{
		requestFinder:    requestFinder,
		roomshipFinder:   roomshipFinder,
		updater:          updater,
		eventIDGenerator: eventIDGenerator,
	}, nil
}

func (uc *agreeMemberRequestUseCase) Execute(ctx context.Context, input *AgreeMemberRequestInput) (*kernel.NoOutput, error) {
	req, err := uc.requestFinder.FindByID(ctx, input.RequestID)
	if err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return nil, myErrors.WrapBusiness(err, "member request not found")
		}
		return nil, err
	}

	roomship, err := uc.roomshipFinder.FindByUserIDAndRoomID(ctx, input.UserID, req.RoomID())
	if err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return nil, domain.ErrNotAdmin
		}
		return nil, err
	}

	if roomship == nil || !roomship.IsAdmin() {
		return nil, domain.ErrNotAdmin
	}

	if err := req.Agree(
		input.UserID,
		uc.eventIDGenerator,
	); err != nil {
		return nil, err
	}

	if err := uc.updater.Update(ctx, req); err != nil {
		return nil, err
	}

	return nil, nil
}
