package application

import (
	"context"
	"encoding/json"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type NotificationIntentCreatedUseCase kernel.UseCase[*NotificationIntentCreatedInput, *kernel.NoOutput]

type NotificationIntentCreatedInput struct {
	ID          kernel.MessageID
	RecipientID kernel.UserID
	RawPayload  json.RawMessage
}

func (r *NotificationIntentCreatedInput) Validate() error {
	if len(r.ID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "notification id can't be empty")
	}
	
	if len(r.RecipientID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "recipient id can't be empty")
	}
	
	return nil
}

type motificationIntentCreatedUseCase struct {
	delivery domain.DeliveryService
}

func NewNotificationIntentCreatedUseCase(
	delivery domain.DeliveryService,
) (NotificationIntentCreatedUseCase, error) {
	if err := validate.NotNil(
		delivery,
	); err != nil {
		return nil, err
	}
	
	return &motificationIntentCreatedUseCase{
		delivery: delivery,
	}, nil
}

func (uc *motificationIntentCreatedUseCase) Execute(ctx context.Context, input *NotificationIntentCreatedInput) (*kernel.NoOutput, error) {
	_ = uc.delivery.Deliver(ctx, input.RecipientID, domain.ActionPushNotification, input.RawPayload)
	return nil, nil
}
