package application

import (
	"context"
	"gochat/internal/chat/domain"
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
	roomMemberDeleter domain.RoomMemberDeleter
}

func NewRoomLeftUseCase(
	roomMemberDeleter domain.RoomMemberDeleter,
) RoomLeftUseCase {
	return &roomLeftUseCase{
		roomMemberDeleter: roomMemberDeleter,
	}
}

func (uc *roomLeftUseCase) Execute(ctx context.Context, input *RoomLeftInput) (*kernel.NoOutput, error) {
	return nil, uc.roomMemberDeleter.DeleteMember(ctx, input.RoomID, input.UserID)
}
