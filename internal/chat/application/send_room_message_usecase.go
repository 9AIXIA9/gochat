package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
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
		return myErrors.ErrEmptyInput
	}
	if len(i.SenderID) == 0 || len(i.RoomID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type sendRoomMessageUseCase struct {
	messageIDGenerator kernel.MessageIDGenerator
	notifier           domain.RoomMessageNotifier
	finder             domain.RoomshipsFinderByRoomID
	messageCreator     domain.RoomMessageCreator
}

func NewSendRoomMessageUseCase(
	messageIDGenerator kernel.MessageIDGenerator,
	notifier domain.RoomMessageNotifier,
	finder domain.RoomshipsFinderByRoomID,
	messageCreator domain.RoomMessageCreator,
) (SendRoomMessageUseCase, error) {
	if err := utils.CheckInterfaces(
		messageIDGenerator,
		notifier,
		finder,
		messageCreator,
	); err != nil {
		return nil, err
	}

	return &sendRoomMessageUseCase{
		messageIDGenerator: messageIDGenerator,
		notifier:           notifier,
		finder:             finder,
		messageCreator:     messageCreator,
	}, nil
}

func (uc *sendRoomMessageUseCase) Execute(ctx context.Context, input *SendRoomMessageInput) (*kernel.NoOutput, error) {
	roomships, err := uc.finder.FindsByRoomID(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	var exist bool
	recipientIDs := make([]kernel.UserID, 0, len(roomships))
	for _, roomship := range roomships {
		if roomship.UserID() == input.SenderID {
			exist = true
			break
		}
		recipientIDs = append(recipientIDs, roomship.UserID())
	}

	if !exist {
		return nil, domain.ErrNotMember
	}

	message, err := domain.CreateRoomMessage(
		input.RoomID,
		input.SenderID,
		recipientIDs,
		input.Content,
		uc.messageIDGenerator,
		uc.notifier,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.messageCreator.Create(ctx, message); err != nil {
		return nil, err
	}

	return nil, nil
}
