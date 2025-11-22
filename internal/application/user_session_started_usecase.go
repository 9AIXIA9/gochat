package application

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
	creator     event.UnpublishedEventCreator
}

func NewUserSessionStartedUseCase(
	idGenerator event.IDGenerator,
	creator event.UnpublishedEventCreator,
) UserSessionStartedUseCase {
	return &userCreatedUseCase{
		idGenerator: idGenerator,
		creator:     creator,
	}
}

func (uc *userCreatedUseCase) Execute(ctx context.Context, input *UserSessionStartedInput) (*kernel.NoOutput, error) {
	ev, err := notificationDomain.NewUndeliveredMessagesNotificationRequestedEvent(uc.idGenerator.Generate(), input.UserID)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.CreateUnpublishedEvent(ctx, ev); err != nil {
		return nil, err
	}
	return nil, nil
}
