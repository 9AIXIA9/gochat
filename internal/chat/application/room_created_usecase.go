package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type RoomCreatedUseCase kernel.UseCase[*RoomCreatedInput, *kernel.NoOutput]

type RoomCreatedInput struct {
	OwnerID    kernel.UserID
	RoomID     kernel.RoomID
	RoomNumber kernel.RoomNumber
}

func (r *RoomCreatedInput) Validate() error {
	if len(r.RoomID) == 0 || len(r.OwnerID) == 0 {
		return myErrors.ErrEmptyInput
	}

	if err := r.RoomNumber.Validate(); err != nil {
		return err
	}
	return nil
}

type roomCreatedUseCase struct {
	roomNumberSaver domain.RoomNumberSaver
	roomMemberSaver domain.RoomMemberSaver
	unitOfWork      kernel.UnitOfWork
}

func NewRoomCreatedUseCase(
	roomNumberSaver domain.RoomNumberSaver,
	roomMemberSaver domain.RoomMemberSaver,
	unitOfWork kernel.UnitOfWork,
) RoomCreatedUseCase {
	return &roomCreatedUseCase{
		roomNumberSaver: roomNumberSaver,
		roomMemberSaver: roomMemberSaver,
		unitOfWork:      unitOfWork,
	}
}

func (uc *roomCreatedUseCase) Execute(ctx context.Context, input *RoomCreatedInput) (*kernel.NoOutput, error) {
	if err := uc.unitOfWork.Execute(ctx, func(txCtx context.Context) error {
		if err := uc.roomNumberSaver.SaveNumber(txCtx, input.RoomID, input.RoomNumber); err != nil {
			return err
		}

		if err := uc.roomMemberSaver.SaveMember(txCtx, input.RoomID, input.OwnerID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return nil, nil
}
