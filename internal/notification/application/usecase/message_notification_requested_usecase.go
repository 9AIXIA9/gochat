package usecase

import (
	"context"
	"errors"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"time"
)

type MessageNotificationRequestedUseCase kernel.UseCase[*MessageNotificationRequestedInput, *kernel.NoOutput]

type MessageNotificationRequestedInput struct {
	MessageID   domain.MessageID
	RecipientID kernel.UserID
	SenderID    kernel.UserID
	Content     string
	SentAt      time.Time
}

func (r *MessageNotificationRequestedInput) Validate() error {
	if len(r.RecipientID) == 0 || len(r.SenderID) == 0 || len(r.MessageID) == 0 {
		return myErrors.ErrEmptyInput
	}

	if len(r.Content) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type messageNotificationRequestedUseCase struct {
	messageSaver         application.MessageSaver
	messageNotifier      application.MessageNotifier
	messageStatesUpdater application.MessageStatesUpdater
}

func NewMessageNotificationRequestedUseCase(
	messageSaver application.MessageSaver,
	messageNotifier application.MessageNotifier,
	messageStatesUpdater application.MessageStatesUpdater,
) MessageNotificationRequestedUseCase {
	return &messageNotificationRequestedUseCase{
		messageSaver:         messageSaver,
		messageNotifier:      messageNotifier,
		messageStatesUpdater: messageStatesUpdater,
	}
}

func (uc *messageNotificationRequestedUseCase) Execute(ctx context.Context, input *MessageNotificationRequestedInput) (*kernel.NoOutput, error) {
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

	if err := uc.messageNotifier.Notify(input.RecipientID, message); err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if err := uc.messageStatesUpdater.UpdateMessageStates(ctx, input.RecipientID, []domain.MessageID{input.MessageID}, domain.MessageStateDelivered); err != nil {
		return nil, err
	}

	return nil, nil
}
