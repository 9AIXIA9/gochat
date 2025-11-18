package usecase

import (
	"context"
	"gochat/internal/chat/application"
	notificationDomain "gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

var _ PrivateMessageCreatedUseCase = (*privateMessageCreatedUseCase)(nil)

type PrivateMessageCreatedUseCase kernel.UseCase[*PrivateMessageCreatedInput, *kernel.NoOutput]

type PrivateMessageCreatedInput struct {
	RecipientID kernel.UserID
	MessageID   kernel.MessageID
}

func (i *PrivateMessageCreatedInput) Validate() error {
	if len(i.RecipientID) == 0 || len(i.MessageID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type privateMessageCreatedUseCase struct {
	eventIDGenerator     event.IDGenerator
	saver                event.UnpublishedEventSaver
	privateMessageFinder application.PrivateMessageFinder
}

func NewPrivateMessageCreatedUseCase(
	eventIDGenerator event.IDGenerator,
	privateMessageFinder application.PrivateMessageFinder,
	saver event.UnpublishedEventSaver,
) PrivateMessageCreatedUseCase {
	return &privateMessageCreatedUseCase{
		eventIDGenerator:     eventIDGenerator,
		saver:                saver,
		privateMessageFinder: privateMessageFinder,
	}
}

func (uc *privateMessageCreatedUseCase) Execute(ctx context.Context, input *PrivateMessageCreatedInput) (*kernel.NoOutput, error) {
	message, err := uc.privateMessageFinder.FindPrivateMessage(ctx, input.RecipientID, input.MessageID)
	if err != nil {
		return nil, err
	}

	if input.RecipientID == message.Sender() {
		//发送给自己的消息不发送通知
		return nil, nil
	}

	ev, err := notificationDomain.NewMessageNotificationRequestedEvent(
		uc.eventIDGenerator.Generate(),
		message.ID(),
		input.RecipientID,
		message.Sender(),
		message.Content(),
		message.SentAt(),
	)
	if err != nil {
		return nil, err
	}

	if err := uc.saver.Save(ctx, ev); err != nil {
		return nil, err
	}

	return nil, nil
}
