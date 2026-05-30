package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

var _ SendRoomMessageUseCase = (*sendRoomMessageUseCase)(nil)

type SendRoomMessageUseCase kernel.UseCase[*SendRoomMessageInput, *kernel.NoOutput]

type SendRoomMessageInput struct {
	SenderID kernel.UserID
	RoomID   kernel.RoomID
	Content  string
}

func (i *SendRoomMessageInput) Validate() error {
	if len(i.Content) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "content can't be empty")
	}
	if len(i.SenderID) == 0 || len(i.RoomID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id and room id can't be empty")
	}
	return nil
}

type sendRoomMessageUseCase struct {
	messageIDGenerator kernel.MessageIDGenerator
	eventIDGenerator   event.IDGenerator
	finder             domain.RoomshipsFinderByRoomID
	messageCreator     domain.RoomMessageCreator
}

func NewSendRoomMessageUseCase(
	messageIDGenerator kernel.MessageIDGenerator,
	eventIDGenerator event.IDGenerator,
	finder domain.RoomshipsFinderByRoomID,
	messageCreator domain.RoomMessageCreator,
) (SendRoomMessageUseCase, error) {
	if err := validate.NotNil(
		messageIDGenerator,
		eventIDGenerator,
		finder,
		messageCreator,
	); err != nil {
		return nil, err
	}

	return &sendRoomMessageUseCase{
		messageIDGenerator: messageIDGenerator,
		finder:             finder,
		messageCreator:     messageCreator,
		eventIDGenerator:   eventIDGenerator,
	}, nil
}

func (uc *sendRoomMessageUseCase) Execute(ctx context.Context, input *SendRoomMessageInput) (*kernel.NoOutput, error) {
	message, err := domain.CreateRoomMessage(
		ctx,
		input.RoomID,
		input.SenderID,
		input.Content,
		uc.finder,
		uc.messageIDGenerator,
		uc.eventIDGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.messageCreator.Create(ctx, message); err != nil {
		return nil, err
	}

	return nil, nil
}
