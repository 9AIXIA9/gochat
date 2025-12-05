package application

import (
	"context"
	"errors"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"time"
)

type PrivateMessageNotificationRequestedUseCase kernel.UseCase[*PrivateMessageNotificationRequestedInput, *kernel.NoOutput]

type PrivateMessageNotificationRequestedInput struct {
	MessageID   kernel.MessageID
	RecipientID kernel.UserID
	SenderID    kernel.UserID
	Content     string
	SentAt      time.Time
}

func (r *PrivateMessageNotificationRequestedInput) Validate() error {
	if len(r.RecipientID) == 0 || len(r.SenderID) == 0 || len(r.MessageID) == 0 {
		return myErrors.ErrEmptyInput
	}

	if len(r.Content) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type privateMessageNotificationRequestedUseCase struct {
	messageSaver    domain.PrivateMessageSaver
	messageNotifier domain.PrivateMessageNotifier
}

func NewPrivateMessageNotificationRequestedUseCase(
	messageSaver domain.PrivateMessageSaver,
	messageNotifier domain.PrivateMessageNotifier,
) PrivateMessageNotificationRequestedUseCase {
	return &privateMessageNotificationRequestedUseCase{
		messageSaver:    messageSaver,
		messageNotifier: messageNotifier,
	}
}

func (uc *privateMessageNotificationRequestedUseCase) Execute(ctx context.Context, input *PrivateMessageNotificationRequestedInput) (*kernel.NoOutput, error) {
	message := domain.LoadPrivateMessage(
		input.MessageID,
		input.SenderID,
		input.RecipientID,
		domain.MessageStateUndelivered,
		input.Content,
		input.SentAt,
	)

	if err := message.Deliver(
		uc.messageNotifier,
	); err != nil && !errors.Is(err, myErrors.ErrNotFound) {
		return nil, err
	}

	if err := uc.messageSaver.Save(ctx, message); err != nil {
		return nil, err
	}
	return nil, nil
}
