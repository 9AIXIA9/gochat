package usecase

import (
	"context"
	"gochat/internal/chat/application"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type RoomLeftUseCase kernel.UseCase[*RoomLeftInput, *kernel.NoOutput]

type RoomLeftInput struct {
	RoomID kernel.RoomID
	UserID kernel.UserID
}

func (r *RoomLeftInput) Validate() error {
	if len(r.RoomID) == 0 || len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type roomLeftUseCase struct {
	roomMemberDeleter application.RoomMemberDeleter
}

func NewRoomLeftUseCase(
	roomMemberDeleter application.RoomMemberDeleter,
) RoomLeftUseCase {
	return &roomLeftUseCase{
		roomMemberDeleter: roomMemberDeleter,
	}
}

func (uc *roomLeftUseCase) Execute(ctx context.Context, input *RoomLeftInput) (*kernel.NoOutput, error) {
	return nil, uc.roomMemberDeleter.DeleteMember(ctx, input.RoomID, input.UserID)
}
