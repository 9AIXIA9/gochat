package usecase

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

var _ SendPrivateMessageUseCase = (*sendPrivateMessageUseCase)(nil)

type SendPrivateMessageUseCase kernel.UseCase[*SendPrivateMessageInput, *kernel.NoOutput]

type SendPrivateMessageInput struct {
	SenderID    kernel.UserID
	RecipientID kernel.UserID
	Content     string
}

func (i *SendPrivateMessageInput) Validate() error {
	if len(i.Content) == 0 {
		return myErrors.ErrEmptyInput
	}
	if len(i.SenderID) == 0 || len(i.RecipientID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type sendPrivateMessageUseCase struct {
	messageIDGenerator    application.MessageIDGenerator
	eventIDGenerator      event.IDGenerator
	userExister           application.UserExister
	messageSaver          application.PrivateMessagesSaver
	unpublishedEventSaver event.UnpublishedSaver
}

func NewSendPrivateMessageUseCase(
	messageIDGenerator application.MessageIDGenerator,
	eventIDGenerator event.IDGenerator,
	userExister application.UserExister,
	messageSaver application.PrivateMessagesSaver,
	unpublishedEventSaver event.UnpublishedSaver,
) SendPrivateMessageUseCase {
	return &sendPrivateMessageUseCase{
		messageIDGenerator:    messageIDGenerator,
		eventIDGenerator:      eventIDGenerator,
		userExister:           userExister,
		messageSaver:          messageSaver,
		unpublishedEventSaver: unpublishedEventSaver,
	}
}

func (uc *sendPrivateMessageUseCase) Execute(ctx context.Context, input *SendPrivateMessageInput) (*kernel.NoOutput, error) {
	if ok, err := uc.userExister.ExistsByID(ctx, input.RecipientID); err != nil {
		return nil, err
	} else if !ok {
		return nil, myErrors.ErrNotFound
	}

	message := domain.NewPrivateMessage(
		uc.messageIDGenerator.Generate(),
		input.SenderID,
		input.Content,
		time.Now().UTC(),
		domain.MessageStateCreated,
		input.RecipientID,
	)

	user := domain.NewUser(input.SenderID, make([]*domain.PrivateMessage, 0, 1))

	if err := user.SendPrivateMessage(message, uc.eventIDGenerator); err != nil {
		return nil, err
	}

	if err := uc.messageSaver.SavePrivateMessages(ctx, user.PrivateMessages()); err != nil {
		return nil, err
	}

	events := user.GetEvents()
	if err := uc.unpublishedEventSaver.Saves(ctx, events); err != nil {
		return nil, err
	}
	return nil, nil
}
