package usecase

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"time"
)

type PrivateMessageCreatedUseCase kernel.UseCase[*PrivateMessageCreatedInput, *kernel.NoOutput]

type PrivateMessageCreatedInput struct {
	MessageID   domain.MessageID
	RecipientID kernel.UserID
	SenderID    kernel.UserID
	Content     string
	SentAt      time.Time
}

func (r *PrivateMessageCreatedInput) Validate() error {
	if len(r.RecipientID) == 0 || len(r.SenderID) == 0 || len(r.MessageID) == 0 {
		return myErrors.ErrEmptyInput
	}

	if len(r.Content) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type privateMessageCreatedUseCase struct {
	messageSaver    application.MessageSaver
	messageNotifier application.MessageNotifier
}

func NewPrivateMessageCreatedUseCase(
	messageSaver application.MessageSaver,
	messageNotifier application.MessageNotifier,
) PrivateMessageCreatedUseCase {
	return &privateMessageCreatedUseCase{
		messageSaver:    messageSaver,
		messageNotifier: messageNotifier,
	}
}

func (uc *privateMessageCreatedUseCase) Execute(ctx context.Context, input *PrivateMessageCreatedInput) (*kernel.NoOutput, error) {
	message := domain.NewMessage(
		input.MessageID,
		input.SenderID,
		domain.MessageStateReceived,
		input.Content,
		input.SentAt,
	)

	if err := uc.messageSaver.SaveMessage(ctx, input.RecipientID, message); err != nil {
		return nil, err
	}

	if err := uc.messageNotifier.Notify(ctx, input.RecipientID, message); err != nil {
		return nil, err
	}

	return nil, nil
}
