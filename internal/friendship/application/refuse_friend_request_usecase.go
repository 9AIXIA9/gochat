package application

import (
	"context"
	"gochat/internal/friendship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
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
		return myErrors.ErrEmptyInput
	}

	return nil
}

type refuseFriendRequestUseCase struct {
	finder      domain.UserFinderByID
	saver       domain.UserSaver
	idGenerator event.IDGenerator
}

func NewRefuseFriendRequestUseCase(
	finder domain.UserFinderByID,
	saver domain.UserSaver,
	idGenerator event.IDGenerator,
) (RefuseFriendRequestUseCase, error) {
	if err := utils.CheckInterfaces(finder, saver, idGenerator); err != nil {
		return nil, err
	}

	return &refuseFriendRequestUseCase{
		finder:      finder,
		saver:       saver,
		idGenerator: idGenerator,
	}, nil
}

func (uc *refuseFriendRequestUseCase) Execute(ctx context.Context, input *RefuseFriendRequestInput) (*kernel.NoOutput, error) {
	user, err := uc.finder.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	if err := user.RefuseFriendRequest(input.RequestID, uc.idGenerator); err != nil {
		return nil, err
	}

	if err := uc.saver.Save(ctx, user); err != nil {
		return nil, err
	}

	return nil, nil
}
