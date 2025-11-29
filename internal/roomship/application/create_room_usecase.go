package application

import (
	"context"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const defaultMemberCount = 20

type CreateRoomUseCase kernel.UseCase[*CreateRoomInput, *CreateRoomOutput]

type CreateRoomInput struct {
	OwnerID        kernel.UserID
	MaxMemberCount int
	Password       domain.Password
}

func (r *CreateRoomInput) Validate() error {
	if err := r.Password.Validate(); err != nil {
		return err
	}
	if len(r.OwnerID) == 0 {
		return myErrors.ErrEmptyInput
	}
	if r.MaxMemberCount < 2 {
		r.MaxMemberCount = defaultMemberCount
	}
	return nil
}

type CreateRoomOutput struct {
	RoomNumber kernel.RoomNumber
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
) CreateRoomUseCase {
	return &createRoomUseCase{
		eventIDGenerator:    eventIDGenerator,
		roomIDGenerator:     roomIDGenerator,
		roomNumberGenerator: numberGenerator,
		encryptor:           encryptor,
		roomCreator:         roomCreator,
	}
}

func (uc *createRoomUseCase) Execute(ctx context.Context, input *CreateRoomInput) (*CreateRoomOutput, error) {
	passwordEncrypted, err := input.Password.Encrypt(uc.encryptor)
	if err != nil {
		return nil, err
	}

	room, err := domain.CreateRoom(
		input.OwnerID,
		uc.roomIDGenerator,
		uc.roomNumberGenerator,
		uc.eventIDGenerator,
		&domain.RoomOption{
			MaxMemberCount:    input.MaxMemberCount,
			PasswordEncrypted: passwordEncrypted,
		},
	)

	if err != nil {
		return nil, err
	}

	if err := uc.roomCreator.Create(ctx, room); err != nil {
		return nil, err
	}

	return &CreateRoomOutput{
		RoomNumber: room.Number(),
	}, nil
}
