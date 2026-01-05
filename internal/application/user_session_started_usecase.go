package application

import (
	"context"
	chatDomain "gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type UserSessionStartedUseCase kernel.UseCase[*UserSessionStartedInput, *kernel.NoOutput]

type UserSessionStartedInput struct {
	UserID kernel.UserID
}

func (r *UserSessionStartedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "UserID cannot be empty")
	}
	return nil
}

type userCreatedUseCase struct {
	idGenerator event.IDGenerator
	creator     event.UnpublishedEventsCreator
}

func NewUserSessionStartedUseCase(
	idGenerator event.IDGenerator,
	creator event.UnpublishedEventsCreator,
) (UserSessionStartedUseCase, error) {
	if err := validate.NotNil(
		idGenerator,
		creator,
	); err != nil {
		return nil, err
	}

	return &userCreatedUseCase{
		idGenerator: idGenerator,
		creator:     creator,
	}, nil
}

func (uc *userCreatedUseCase) Execute(ctx context.Context, input *UserSessionStartedInput) (*kernel.NoOutput, error) {
	chatEvPushRequestedEvent, err := chatDomain.NewUndeliveredMessagesPushRequestedEvent(input.UserID, uc.idGenerator)
	if err != nil {
		return nil, err
	}

	notificationEvPushRequestedEvent, err := notificationDomain.NewUndeliveredMessagesNotificationRequestedEvent(input.UserID, uc.idGenerator)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.CreateUnpublishedEvents(ctx, []event.Event{chatEvPushRequestedEvent, notificationEvPushRequestedEvent}); err != nil {
		return nil, err
	}
	return nil, nil
}
