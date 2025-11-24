package application

import (
	"context"
	"gochat/internal/authorization/domain"
	chatDomain "gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	socialDomain "gochat/internal/social/domain"
)

type UserCreatedUseCase kernel.UseCase[*UserCreatedInput, *kernel.NoOutput]

type UserCreatedInput struct {
	UserID kernel.UserID
}

func (r *UserCreatedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type userCreatedUseCase struct {
	idGenerator event.IDGenerator
	creator     event.UnpublishedEventsCreator
	finder      domain.UserFinderByID
}

func NewUserCreatedUseCase(
	idGenerator event.IDGenerator,
	creator event.UnpublishedEventsCreator,
	finder domain.UserFinderByID,
) UserCreatedUseCase {
	return &userCreatedUseCase{
		idGenerator: idGenerator,
		creator:     creator,
		finder:      finder,
	}
}

func (uc *userCreatedUseCase) Execute(ctx context.Context, input *UserCreatedInput) (*kernel.NoOutput, error) {
	user, err := uc.finder.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	socialEv, err := socialDomain.NewUserCreatedEvent(user.ID(), uc.idGenerator)
	if err != nil {
		return nil, err
	}

	chatEv, err := chatDomain.NewUserCreatedEvent(
		user.ID(),
		user.Number(),
		uc.idGenerator,
	)
	if err != nil {
		return nil, err
	}

	notificationEv, err := notificationDomain.NewWelcomeEmailNotificationRequestedEvent(
		user.ID(),
		user.Email(),
		user.Number(),
		uc.idGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.CreateUnpublishedEvents(ctx, []event.Event{
		socialEv, chatEv, notificationEv,
	}); err != nil {
		return nil, err
	}
	return nil, nil
}
