package application

import (
	"context"
	chatDomain "gochat/internal/chat/domain"
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
	chatEvPushRequestedEvent, err := chatDomain.NewUndeliveredMessagesPushRequestedEvent(input.UserID, uc.idGenerator)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.CreateUnpublishedEvent(ctx, chatEvPushRequestedEvent); err != nil {
		return nil, err
	}
	return nil, nil
}
