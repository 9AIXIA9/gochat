package application

import (
	"context"
	"gochat/internal/profile/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
	"time"
)

type RoomCreatedUseCase kernel.UseCase[*RoomCreatedInput, *kernel.NoOutput]

type RoomCreatedInput struct {
	RoomID    kernel.RoomID
	CreatedAt time.Time
}

func (r *RoomCreatedInput) Validate() error {
	if len(r.RoomID) == 0 || r.CreatedAt.IsZero() {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "room id or created time is empty")
	}

	return nil
}

type roomCreatedUseCase struct {
	roomSaver domain.RoomSaver
	creator   domain.RoomProfileCreator
}

func NewRoomCreatedUseCase(
	roomSaver domain.RoomSaver,
	creator domain.RoomProfileCreator,
) (RoomCreatedUseCase, error) {
	if err := utils.CheckInterfaces(
		roomSaver,
		creator,
	); err != nil {
		return nil, err
	}

	return &roomCreatedUseCase{
		roomSaver: roomSaver,
		creator:   creator,
	}, nil
}

func (uc *roomCreatedUseCase) Execute(ctx context.Context, input *RoomCreatedInput) (*kernel.NoOutput, error) {
	if err := uc.roomSaver.Save(ctx, domain.LoadRoom(input.RoomID)); err != nil {
		return nil, err
	}

	profile := domain.CreateRoomProfile(
		input.RoomID,
		input.CreatedAt,
	)

	if err := uc.creator.Create(ctx, profile); err != nil {
		return nil, err
	}
	return nil, nil
}
