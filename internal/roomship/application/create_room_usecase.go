package application

import (
	"context"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

const defaultMemberCount = 20

type CreateRoomUseCase kernel.UseCase[*CreateRoomInput, *kernel.NoOutput]

type CreateRoomInput struct {
	UserID         kernel.UserID
	MaxMemberCount int
	Password       domain.Password
}

func (r *CreateRoomInput) Validate() error {
	if err := r.Password.Validate(); err != nil {
		return err
	}
	if len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id is empty")
	}
	if r.MaxMemberCount < 2 {
		r.MaxMemberCount = defaultMemberCount
	}
	return nil
}

type createRoomUseCase struct {
	eventIDGenerator    event.IDGenerator
	roomIDGenerator     domain.RoomIDGenerator
	roomNumberGenerator domain.RoomNumberGenerator
	encryptor           domain.Encryptor
	roomCreator         domain.RoomCreator
}

func NewCreateRoomUseCase(
	eventIDGenerator event.IDGenerator,
	roomIDGenerator domain.RoomIDGenerator,
	numberGenerator domain.RoomNumberGenerator,
	encryptor domain.Encryptor,
	roomCreator domain.RoomCreator,
) (CreateRoomUseCase, error) {
	if err := validate.NotNil(
		eventIDGenerator,
		roomIDGenerator,
		numberGenerator,
		encryptor,
		roomCreator,
	); err != nil {
		return nil, err
	}
	return &createRoomUseCase{
		eventIDGenerator:    eventIDGenerator,
		roomIDGenerator:     roomIDGenerator,
		roomNumberGenerator: numberGenerator,
		encryptor:           encryptor,
		roomCreator:         roomCreator,
	}, nil
}

func (uc *createRoomUseCase) Execute(ctx context.Context, input *CreateRoomInput) (*kernel.NoOutput, error) {
	passwordEncrypted, err := input.Password.Encrypt(uc.encryptor)
	if err != nil {
		return nil, err
	}

	room, err := domain.CreateRoom(
		input.UserID,
		input.MaxMemberCount,
		passwordEncrypted,
		uc.roomIDGenerator,
		uc.roomNumberGenerator,
		uc.eventIDGenerator,
	)

	if err != nil {
		return nil, err
	}

	if err := uc.roomCreator.Create(ctx, room); err != nil {
		return nil, err
	}

	return nil, nil
}
