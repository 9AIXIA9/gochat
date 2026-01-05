package application

import (
	"context"
	"gochat/internal/authorization/domain"
	chatDomain "gochat/internal/chat/domain"
	friendshipDomain "gochat/internal/friendship/domain"
	notificationDomain "gochat/internal/notification/domain"
	profileDomain "gochat/internal/profile/domain"
	roomshipDomain "gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type UserCreatedUseCase kernel.UseCase[*UserCreatedInput, *kernel.NoOutput]

type UserCreatedInput struct {
	UserID kernel.UserID
}

func (r *UserCreatedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "UserID cannot be empty")
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
) (UserCreatedUseCase, error) {
	if err := validate.NotNil(
		idGenerator,
		creator,
		finder,
	); err != nil {
		return nil, err
	}

	return &userCreatedUseCase{
		idGenerator: idGenerator,
		creator:     creator,
		finder:      finder,
	}, nil
}

func (uc *userCreatedUseCase) Execute(ctx context.Context, input *UserCreatedInput) (*kernel.NoOutput, error) {
	user, err := uc.finder.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	profileEv, err := profileDomain.NewUserCreatedEvent(
		user.ID(),
		user.Email(),
		user.SignedUpAt(),
		uc.idGenerator,
	)
	if err != nil {
		return nil, err
	}

	friendshipEv, err := friendshipDomain.NewUserCreatedEvent(user.ID(), uc.idGenerator)
	if err != nil {
		return nil, err
	}

	roomshipEv, err := roomshipDomain.NewUserCreatedEvent(user.ID(), uc.idGenerator)
	if err != nil {
		return nil, err
	}

	chatEv, err := chatDomain.NewUserCreatedEvent(
		user.ID(),
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
		profileEv, friendshipEv, roomshipEv, chatEv, notificationEv,
	}); err != nil {
		return nil, err
	}
	return nil, nil
}
