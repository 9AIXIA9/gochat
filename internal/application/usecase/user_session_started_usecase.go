package usecase

import (
	"context"
	notificationDomain "gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type UserSessionStartedUseCase kernel.UseCase[*UserSessionStartedInput, *kernel.NoOutput]

type UserSessionStartedInput struct {
	UserID kernel.UserID
}

func (r *UserSessionStartedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type userCreatedUseCase struct {
	idGenerator event.IDGenerator
	saver       event.UnpublishedEventSaver
}

func NewUserSessionStartedUseCase(
	idGenerator event.IDGenerator,
	saver event.UnpublishedEventSaver,
) UserSessionStartedUseCase {
	return &userCreatedUseCase{
		idGenerator: idGenerator,
		saver:       saver,
	}
}

func (uc *userCreatedUseCase) Execute(ctx context.Context, input *UserSessionStartedInput) (*kernel.NoOutput, error) {
	notificationEv, err := notificationDomain.NewUndeliveredMessageNotificationRequestedEvent(
		uc.idGenerator.Generate(),
		input.UserID,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.saver.Save(ctx, notificationEv); err != nil {
		return nil, err
	}
	return nil, nil
}
