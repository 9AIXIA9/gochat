package usecase

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

var _ PrivateMessageCreatedUseCase = (*privateMessageCreatedUseCase)(nil)

type PrivateMessageCreatedUseCase kernel.UseCase[*PrivateMessageCreatedInput, *kernel.NoOutput]

type PrivateMessageCreatedInput struct {
	RecipientID kernel.UserID
	MessageID   domain.MessageID
}

func (i *PrivateMessageCreatedInput) Validate() error {
	if len(i.RecipientID) == 0 || len(i.MessageID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type privateMessageCreatedUseCase struct {
	eventIDGenerator     event.IDGenerator
	publisher            event.Publisher
	privateMessageFinder application.PrivateMessageFinder
}

func NewPrivateMessageCreatedUseCase(
	eventIDGenerator event.IDGenerator,
	privateMessageFinder application.PrivateMessageFinder,
	publisher event.Publisher,
) PrivateMessageCreatedUseCase {
	return &privateMessageCreatedUseCase{
		eventIDGenerator:     eventIDGenerator,
		publisher:            publisher,
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

	ev, err := notificationDomain.NewPrivateMessageCreatedEvent(
		uc.eventIDGenerator.Generate(),
		notificationDomain.MessageID(message.ID()),
		input.RecipientID,
		message.Sender(),
		message.Content(),
		message.SentAt(),
	)
	if err != nil {
		return nil, err
	}

	if err := uc.publisher.Publish(ev); err != nil {
		return nil, err
	}

	return nil, nil
}
