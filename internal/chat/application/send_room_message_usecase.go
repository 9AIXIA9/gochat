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
		return myErrors.NewBusiness("message content cannot be empty")
	}
	if len(i.SenderID) == 0 || len(i.RoomID) == 0 {
		return myErrors.NewBusiness("invalid sender or room ID")
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

	if len(roomships) == 0 {
		return nil, myErrors.NewBusiness("room not found")
	}

	message, err := uc.createRoomMessage(roomships, input)
	if err != nil {
		return nil, err
	}

	if err := uc.messageCreator.Create(ctx, message); err != nil {
		return nil, err
	}

	return nil, nil
}

func (uc *sendRoomMessageUseCase) createRoomMessage(roomships []*domain.Roomship, input *SendRoomMessageInput) (*domain.RoomMessage, error) {
	// 房间内没有其他成员
	if len(roomships) == 1 {
		// 仅有自己
		if roomships[0].UserID() == input.SenderID {
			return domain.CreateRoomMessage(
				input.RoomID,
				input.SenderID,
				nil,
				input.Content,
				uc.messageIDGenerator,
				uc.notifier,
			), nil
		}
		return nil, myErrors.NewBusiness("you are not a member of the room")
	}

	var exist bool
	recipientIDs := make([]kernel.UserID, 0, len(roomships)-1)
	for _, roomship := range roomships {
		if roomship.UserID() == input.SenderID {
			exist = true
			continue
		}
		recipientIDs = append(recipientIDs, roomship.UserID())
	}

	if !exist {
		return nil, myErrors.NewBusiness("you are not a member of the room")
	}

	return domain.CreateRoomMessage(
		input.RoomID,
		input.SenderID,
		recipientIDs,
		input.Content,
		uc.messageIDGenerator,
		uc.notifier,
	), nil
}
