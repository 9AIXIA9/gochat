package usecase

import (
	"context"
	"gochat/internal/authorization/application"
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
	publisher   event.Publisher
	finder      application.UserFinderByID
}

func NewUserCreatedUseCase(
	idGenerator event.IDGenerator,
	publisher event.Publisher,
	finder application.UserFinderByID,
) UserCreatedUseCase {
	return &userCreatedUseCase{
		idGenerator: idGenerator,
		publisher:   publisher,
		finder:      finder,
	}
}

func (uc *userCreatedUseCase) Execute(ctx context.Context, input *UserCreatedInput) (*kernel.NoOutput, error) {
	user, err := uc.finder.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	socialEv, err := socialDomain.NewUserCreatedEvent(uc.idGenerator.Generate(), user.ID())
	if err != nil {
		return nil, err
	}

	if err := uc.publisher.Publish(socialEv); err != nil {
		return nil, err
	}

	chatEv, err := chatDomain.NewUserCreatedEvent(
		uc.idGenerator.Generate(),
		user.ID(),
		chatDomain.UserNumber(user.Number()),
	)
	if err != nil {
		return nil, err
	}

	if err := uc.publisher.Publish(chatEv); err != nil {
		return nil, err
	}

	notificationEv, err := notificationDomain.NewWelcomeEmailNotificationRequestedEvent(
		uc.idGenerator.Generate(),
		user.ID(),
		user.Email(),
		notificationDomain.UserNumber(user.Number()),
	)
	if err != nil {
		return nil, err
	}

	if err := uc.publisher.Publish(notificationEv); err != nil {
		return nil, err
	}
	return nil, nil
}
