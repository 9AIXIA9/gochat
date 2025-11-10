package usecase

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"time"
)

type RoomMessageCreatedUseCase kernel.UseCase[*RoomMessageCreatedInput, *kernel.NoOutput]

type RoomMessageCreatedInput struct {
	MessageID    domain.MessageID
	RoomID       domain.RoomID
	RecipientIDs []kernel.UserID
	SenderID     kernel.UserID
	Content      string
	SentAt       time.Time
}

func (r *RoomMessageCreatedInput) Validate() error {
	if len(r.SenderID) == 0 || len(r.MessageID) == 0 || len(r.RoomID) == 0 || len(r.RecipientIDs) == 0 {
		return myErrors.ErrEmptyInput
	}

	if len(r.Content) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type roomMessageCreatedUseCase struct {
	messageSaver    application.MessageSaver
	messageNotifier application.MessageNotifier
}

func NewRoomMessageCreatedUseCase(messageSaver application.MessageSaver, messageNotifier application.MessageNotifier) RoomMessageCreatedUseCase {
	return &roomMessageCreatedUseCase{
		messageSaver:    messageSaver,
		messageNotifier: messageNotifier,
	}
}

func (uc *roomMessageCreatedUseCase) Execute(ctx context.Context, input *RoomMessageCreatedInput) (*kernel.NoOutput, error) {
	message := domain.NewMessage(
		input.MessageID,
		input.SenderID,
		domain.MessageStateReceived,
		input.Content,
		input.SentAt,
	)

	for _, recipientID := range input.RecipientIDs {
		if err := uc.messageSaver.SaveMessage(ctx, recipientID, message); err != nil {
			return nil, err
		}

		if err := uc.messageNotifier.Notify(ctx, recipientID, message); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
