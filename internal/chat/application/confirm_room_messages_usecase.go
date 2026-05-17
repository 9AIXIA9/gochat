package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type ConfirmRoomMessagesUseCase kernel.UseCase[*ConfirmRoomMessagesInput, *kernel.NoOutput]

type ConfirmRoomMessagesInput struct {
	UserID     kernel.UserID
	MessageIDs []kernel.MessageID
}

func (r *ConfirmRoomMessagesInput) Validate() error {
	if len(r.UserID) == 0 || len(r.MessageIDs) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id and message ids can't be empty")
	}

	return nil
}

type confirmRoomMessagesUseCase struct {
	updater domain.RoomMessagesStatesUpdaterByMessageIDs
}

func NewConfirmRoomMessagesUseCase(
	updater domain.RoomMessagesStatesUpdaterByMessageIDs,
) (ConfirmRoomMessagesUseCase, error) {
	if err := validate.NotNil(updater); err != nil {
		return nil, err
	}

	return &confirmRoomMessagesUseCase{
		updater: updater,
	}, nil
}

func (uc *confirmRoomMessagesUseCase) Execute(ctx context.Context, input *ConfirmRoomMessagesInput) (*kernel.NoOutput, error) {
	return nil, uc.updater.UpdatesByMessageIDs(ctx, input.UserID, input.MessageIDs, domain.MessageStateDelivered)
}
