package application

import (
	"context"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type RoomCreatedUseCase kernel.UseCase[*RoomCreatedInput, *kernel.NoOutput]

type RoomCreatedInput struct {
	RoomID kernel.RoomID
}

func (r *RoomCreatedInput) Validate() error {
	if len(r.RoomID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type roomCreatedUseCase struct {
	roomshipIDGenerator domain.RoomshipIDGenerator
	idGenerator         event.IDGenerator
	finder              domain.RoomFinderByID
	creator             domain.RoomshipCreator
}

func NewRoomCreatedUseCase(
	roomshipIDGenerator domain.RoomshipIDGenerator,
	idGenerator event.IDGenerator,
	finder domain.RoomFinderByID,
	creator domain.RoomshipCreator,
) (RoomCreatedUseCase, error) {
	if err := utils.CheckInterfaces(
		roomshipIDGenerator,
		idGenerator,
		finder,
		creator,
	); err != nil {
		return nil, err
	}

	return &roomCreatedUseCase{
		roomshipIDGenerator: roomshipIDGenerator,
		idGenerator:         idGenerator,
		finder:              finder,
		creator:             creator,
	}, nil
}

func (uc *roomCreatedUseCase) Execute(ctx context.Context, input *RoomCreatedInput) (*kernel.NoOutput, error) {
	room, err := uc.finder.FindByID(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	roomship, err := domain.CreateRoomship(
		room.OwnerID(),
		input.RoomID,
		domain.OwnerRole,
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
