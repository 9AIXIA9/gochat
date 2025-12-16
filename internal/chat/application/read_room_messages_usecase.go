package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type ReadRoomMessagesUseCase kernel.UseCase[*ReadRoomMessagesInput, *kernel.NoOutput]

type ReadRoomMessagesInput struct {
	UserID kernel.UserID
	RoomID kernel.RoomID
}

func (r *ReadRoomMessagesInput) Validate() error {
	if len(r.UserID) == 0 || len(r.RoomID) == 0 {
		return myErrors.NewBusiness("user ID and room ID cannot be empty")
	}

	return nil
}

type readRoomMessagesUseCase struct {
	updater domain.RoomMessagesStatesUpdaterByUserIDAndRoomID
}

func NewReadRoomMessagesUseCase(
	updater domain.RoomMessagesStatesUpdaterByUserIDAndRoomID,
) (ReadRoomMessagesUseCase, error) {
	if err := utils.CheckInterfaces(updater); err != nil {
		return nil, err
	}

	return &readRoomMessagesUseCase{
		updater: updater,
	}, nil
}

func (uc *readRoomMessagesUseCase) Execute(ctx context.Context, input *ReadRoomMessagesInput) (*kernel.NoOutput, error) {
	return nil, uc.updater.UpdatesByUserIDAndRoomID(ctx, input.UserID, input.RoomID, domain.MessageStateRead)
}
