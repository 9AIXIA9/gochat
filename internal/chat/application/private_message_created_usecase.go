package application

import (
	"context"
	"gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
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
	creator              event.UnpublishedEventCreator
	privateMessageFinder domain.PrivateMessageFinder
}

func NewPrivateMessageCreatedUseCase(
	eventIDGenerator event.IDGenerator,
	privateMessageFinder domain.PrivateMessageFinder,
	creator event.UnpublishedEventCreator,
) (PrivateMessageCreatedUseCase, error) {
	if err := utils.CheckInterfaces(
		eventIDGenerator,
		privateMessageFinder,
		creator,
	); err != nil {
		return nil, err
	}

	return &privateMessageCreatedUseCase{
		eventIDGenerator:     eventIDGenerator,
		creator:              creator,
		privateMessageFinder: privateMessageFinder,
	}, nil
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

	ev, err := notificationDomain.NewPrivateMessageNotificationRequestedEvent(
		message.ID(),
		message.RecipientID(),
		message.SenderID(),
		message.Content(),
		message.SentAt(),
		uc.eventIDGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.CreateUnpublishedEvent(ctx, ev); err != nil {
		return nil, err
	}

	return nil, nil
}
