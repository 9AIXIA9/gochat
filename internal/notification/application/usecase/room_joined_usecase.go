package usecase

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type RoomJoinedUseCase kernel.UseCase[*RoomJoinedInput, *kernel.NoOutput]

type RoomJoinedInput struct {
	RoomID domain.RoomID
	UserID kernel.UserID
}

func (r *RoomJoinedInput) Validate() error {
	if len(r.RoomID) == 0 || len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type roomJoinedUseCase struct {
	roomMemberSaver application.RoomMemberSaver
}

func NewRoomJoinedUseCase(
	roomMemberSaver application.RoomMemberSaver,
) RoomJoinedUseCase {
	return &roomJoinedUseCase{
		roomMemberSaver: roomMemberSaver,
	}
}

func (uc *roomJoinedUseCase) Execute(ctx context.Context, input *RoomJoinedInput) (*kernel.NoOutput, error) {
	return nil, uc.roomMemberSaver.SaveMember(ctx, input.RoomID, input.UserID)
}
