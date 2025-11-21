package application

import (
	"context"
	"gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

var _ PrivateMessageCreatedUseCase = (*privateMessageCreatedUseCase)(nil)

type PrivateMessageCreatedUseCase kernel.UseCase[*PrivateMessageCreatedInput, *kernel.NoOutput]

type PrivateMessageCreatedInput struct {
	MessageID kernel.MessageID
}

func (i *PrivateMessageCreatedInput) Validate() error {
	if len(i.MessageID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type privateMessageCreatedUseCase struct {
	eventIDGenerator     event.IDGenerator
	saver                event.UnpublishedEventSaver
	privateMessageFinder domain.PrivateMessageFinder
}

func NewPrivateMessageCreatedUseCase(
	eventIDGenerator event.IDGenerator,
	privateMessageFinder domain.PrivateMessageFinder,
	saver event.UnpublishedEventSaver,
) PrivateMessageCreatedUseCase {
	return &privateMessageCreatedUseCase{
		eventIDGenerator:     eventIDGenerator,
		saver:                saver,
		privateMessageFinder: privateMessageFinder,
	}
}

func (uc *privateMessageCreatedUseCase) Execute(ctx context.Context, input *PrivateMessageCreatedInput) (*kernel.NoOutput, error) {
	message, err := uc.privateMessageFinder.FindPrivateMessage(ctx, input.MessageID)
	if err != nil {
		return nil, err
	}

	if message.SenderID() == message.RecipientID() {
		//发送给自己的消息不发送通知
		return nil, nil
	}

	ev, err := notificationDomain.NewMessageNotificationRequestedEvent(
		uc.eventIDGenerator.Generate(),
		message.ID(),
		message.RecipientID(),
		message.SenderID(),
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
