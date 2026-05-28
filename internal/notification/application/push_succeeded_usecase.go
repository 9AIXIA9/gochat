package application

import (
	"context"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type PushSucceededUseCase kernel.UseCase[*PushSucceededInput, *kernel.NoOutput]

type PushSucceededInput struct {
	ID kernel.MessageID
}

func (r *PushSucceededInput) Validate() error {
	if len(r.ID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "notification id can't be empty")
	}
	return nil
}

type pushSucceededUseCase struct {
	updater domain.NotificationStateUpdater
}

func NewPushSucceededUseCase(
	updater domain.NotificationStateUpdater,
) (PushSucceededUseCase, error) {
	if err := validate.NotNil(
		updater,
	); err != nil {
		return nil, err
	}

	return &pushSucceededUseCase{
		updater: updater,
	}, nil
}

func (uc *pushSucceededUseCase) Execute(ctx context.Context, input *PushSucceededInput) (*kernel.NoOutput, error) {
	if err := uc.updater.UpdateStateByID(ctx, input.ID, domain.StateDelivered); err != nil {
		return nil, err
	}
	return nil, nil
}
