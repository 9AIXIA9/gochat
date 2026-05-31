package application

import (
	"context"
	"encoding/json"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
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
	delivery domain.DeliveryService
}

func NewNotificationCreatedUseCase(
	delivery domain.DeliveryService,
) (NotificationCreatedUseCase, error) {
	if err := validate.NotNil(
		delivery,
	); err != nil {
		return nil, err
	}

	return &notificationCreatedUseCase{
		delivery: delivery,
	}, nil
}

func (uc *notificationCreatedUseCase) Execute(ctx context.Context, input *NotificationCreatedInput) (*kernel.NoOutput, error) {
	_ = uc.delivery.Deliver(ctx, input.RecipientID, domain.ActionPushNotification, input.RawPayload)
	return nil, nil
}
