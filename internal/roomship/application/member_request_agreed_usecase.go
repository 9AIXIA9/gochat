package application

import (
	"context"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type MemberRequestAgreedUseCase kernel.UseCase[*MemberRequestAgreedInput, *kernel.NoOutput]

type MemberRequestAgreedInput struct {
	UserID kernel.UserID
	RoomID kernel.RoomID
}

func (r *MemberRequestAgreedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type memberRequestAgreedUseCase struct {
	roomshipIDGenerator domain.RoomshipIDGenerator
	idGenerator         event.IDGenerator
	creator             domain.RoomshipCreator
}

func NewMemberRequestAgreedUseCase(
	roomshipIDGenerator domain.RoomshipIDGenerator,
	idGenerator event.IDGenerator,
	creator domain.RoomshipCreator,
) MemberRequestAgreedUseCase {
	return &memberRequestAgreedUseCase{
		roomshipIDGenerator: roomshipIDGenerator,
		idGenerator:         idGenerator,
		creator:             creator,
	}
}

func (uc *memberRequestAgreedUseCase) Execute(ctx context.Context, input *MemberRequestAgreedInput) (*kernel.NoOutput, error) {
	roomship, err := domain.CreateRoomship(
		input.UserID,
		input.RoomID,
		domain.MemberRole,
		uc.roomshipIDGenerator,
		uc.idGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.Create(ctx, roomship); err != nil {
		return nil, err
	}

	return nil, nil
}
