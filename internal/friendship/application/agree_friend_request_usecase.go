package application

import (
	"context"
	"gochat/internal/friendship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type AgreeFriendRequestUseCase kernel.UseCase[*AgreeFriendRequestInput, *kernel.NoOutput]

type AgreeFriendRequestInput struct {
	UserID    kernel.UserID
	RequestID kernel.OperationID
}

func (r *AgreeFriendRequestInput) Validate() error {
	if len(r.RequestID) == 0 || len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type agreeFriendRequestUseCase struct {
	finder      domain.UserFinderByID
	saver       domain.UserSaver
	idGenerator event.IDGenerator
}

func NewAgreeFriendRequestUseCase(
	finder domain.UserFinderByID,
	saver domain.UserSaver,
	idGenerator event.IDGenerator,
) (AgreeFriendRequestUseCase, error) {
	if err := utils.CheckInterfaces(
		finder, saver, idGenerator,
	); err != nil {
		return nil, err
	}
	return &agreeFriendRequestUseCase{
		finder:      finder,
		saver:       saver,
		idGenerator: idGenerator,
	}, nil
}

func (uc *agreeFriendRequestUseCase) Execute(ctx context.Context, input *AgreeFriendRequestInput) (*kernel.NoOutput, error) {
	user, err := uc.finder.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	if err := user.AgreeFriendRequest(input.RequestID, uc.idGenerator); err != nil {
		return nil, err
	}

	if err := uc.saver.Save(ctx, user); err != nil {
		return nil, err
	}

	return nil, nil
}
