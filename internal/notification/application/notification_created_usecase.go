package application

import (
	"context"
	"encoding/json"
	"gochat/internal/gateway/core"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/contract"
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
	gateway contract.GatewayService
}

func NewNotificationCreatedUseCase(
	gateway contract.GatewayService,
) (NotificationCreatedUseCase, error) {
	if err := validate.NotNil(
		gateway,
	); err != nil {
		return nil, err
	}

	return &notificationCreatedUseCase{
		gateway: gateway,
	}, nil
}

func (uc *notificationCreatedUseCase) Execute(ctx context.Context, input *NotificationCreatedInput) (*kernel.NoOutput, error) {
	_ = uc.gateway.PushToUser(ctx, input.RecipientID, core.NewActiveDownstreamEnvelop(domain.ActionPushNotification, input.RawPayload))
	return nil, nil
}
