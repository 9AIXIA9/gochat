package usecase

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

var _ SendMessageUseCase = (*sendMessageUseCase)(nil)

type SendMessageUseCase kernel.UseCase[*SendMessageInput, *kernel.NoOutput]

type SendMessageInput struct {
	MessageID domain.MessageID
	Sender    kernel.UserID
	Recipient kernel.UserID
	Content   string
	SentAt    time.Time
}

func (i *SendMessageInput) Validate() error {
	if len(i.Sender) == 0 || len(i.Recipient) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type sendMessageUseCase struct {
	eventIDGenerator event.IDGenerator
	messageNotifier  application.MessageNotifier
}

func NewSendMessageUseCase(eventIDGenerator event.IDGenerator, messageNotifier application.MessageNotifier) SendMessageUseCase {
	return &sendMessageUseCase{
		eventIDGenerator: eventIDGenerator,
		messageNotifier:  messageNotifier,
	}
}

func (uc *sendMessageUseCase) Execute(ctx context.Context, input *SendMessageInput) (*kernel.NoOutput, error) {
	message := domain.NewMessage(input.MessageID, input.Recipient, input.Sender, input.Content, input.SentAt)

	if err := uc.messageNotifier.Notify(ctx, message); err != nil {
		return nil, err
	}

	return nil, nil
}
