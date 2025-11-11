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
	messageSaver         application.MessageSaver
	messageNotifier      application.MessageNotifier
	messageStatesUpdater application.MessageStatesUpdater
}

func NewRoomMessageCreatedUseCase(
	messageSaver application.MessageSaver,
	messageNotifier application.MessageNotifier,
	messageStatesUpdater application.MessageStatesUpdater,
) RoomMessageCreatedUseCase {
	return &roomMessageCreatedUseCase{
		messageSaver:         messageSaver,
		messageNotifier:      messageNotifier,
		messageStatesUpdater: messageStatesUpdater,
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

		if err := uc.messageNotifier.Notify(recipientID, message); err != nil {
			if errors.Is(err, myErrors.ErrNotFound) {
				continue
			}
			return nil, err
		}

		if err := uc.messageStatesUpdater.UpdateMessageStates(ctx, recipientID, []domain.MessageID{input.MessageID}, domain.MessageStateDelivered); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
