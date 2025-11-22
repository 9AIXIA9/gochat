package application

import (
	"context"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"
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
	eventIDGenerator event.IDGenerator
	roomIDGenerator  domain.RoomIDGenerator
	numberGenerator  domain.RoomNumberGenerator
	encryptor        domain.Encryptor
	roomSaver        domain.RoomSaver
	eventsCreator    event.UnpublishedEventsCreator
	unitOfWork       kernel.UnitOfWork
}

func NewCreateRoomUseCase(
	eventIDGenerator event.IDGenerator,
	roomIDGenerator domain.RoomIDGenerator,
	numberGenerator domain.RoomNumberGenerator,
	encryptor domain.Encryptor,
	roomSaver domain.RoomSaver,
	eventsCreator event.UnpublishedEventsCreator,
	unitOfWork kernel.UnitOfWork,
) CreateRoomUseCase {
	return &createRoomUseCase{
		eventIDGenerator: eventIDGenerator,
		roomIDGenerator:  roomIDGenerator,
		numberGenerator:  numberGenerator,
		encryptor:        encryptor,
		roomSaver:        roomSaver,
		eventsCreator:    eventsCreator,
		unitOfWork:       unitOfWork,
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
		uc.numberGenerator,
		uc.eventIDGenerator,
		&domain.RoomOption{
			MaxMemberCount:    input.MaxMemberCount,
			PasswordEncrypted: passwordEncrypted,
		},
	)

	if err := uc.unitOfWork.Execute(ctx, func(txCtx context.Context) error {
		if err := uc.roomSaver.Save(txCtx, room); err != nil {
			return err
		}

		if err := uc.eventsCreator.CreateUnpublishedEvents(txCtx, room.GetEvents()); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &CreateRoomOutput{
		RoomNumber: room.Number(),
	}, nil
}
