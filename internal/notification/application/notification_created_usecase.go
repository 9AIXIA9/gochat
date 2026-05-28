package application

import (
	"context"
	"encoding/json"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type NotificationCreatedUseCase kernel.UseCase[*NotificationCreatedInput, *kernel.NoOutput]

type NotificationCreatedInput struct {
	ID          kernel.MessageID
	RecipientID kernel.UserID
	RawPayload  json.RawMessage
}

func (r *NotificationCreatedInput) Validate() error {
	if len(r.ID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "notification id can't be empty")
	}

	if len(r.RecipientID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "recipient id can't be empty")
	}

	return nil
}

type notificationCreatedUseCase struct {
	idGenerator event.IDGenerator
	creator     domain.NotificationCreator
}

func NewNotificationCreatedUseCase(
	idGenerator event.IDGenerator,
	creator domain.NotificationCreator,
) (NotificationCreatedUseCase, error) {
	if err := validate.NotNil(
		idGenerator,
		creator,
	); err != nil {
		return nil, err
	}

	return &notificationCreatedUseCase{
		idGenerator: idGenerator,
		creator:     creator,
	}, nil
}

func (uc *notificationCreatedUseCase) Execute(ctx context.Context, input *NotificationCreatedInput) (*kernel.NoOutput, error) {
	notification, err := domain.CreateNotification(
		input.ID,
		input.RecipientID,
		input.RawPayload,
		uc.idGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.Create(ctx, notification); err != nil {
		return nil, err
	}
	return nil, nil
}
