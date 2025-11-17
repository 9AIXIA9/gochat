package usecase

import (
	"context"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/application"
	"gochat/internal/social/domain"
	"time"
)

const defaultMemberCount = 20

type CreateRoomUseCase kernel.UseCase[*CreateRoomInput, *CreateRoomOutput]

type CreateRoomInput struct {
	Owner          kernel.UserID
	MaxMemberCount int
	Password       domain.Password
}

func (r *CreateRoomInput) Validate() error {
	if err := r.Password.Validate(); err != nil {
		return err
	}
	if len(r.Owner) == 0 {
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
	roomIDGenerator  application.RoomIDGenerator
	numberGenerator  application.RoomNumberGenerator
	encryptor        application.Encryptor
	roomSaver        application.RoomSaver
	eventSaver       event.UnpublishedSaver
	unitOfWork       kernel.UnitOfWork
}

func NewCreateRoomUseCase(
	eventIDGenerator event.IDGenerator,
	roomIDGenerator application.RoomIDGenerator,
	numberGenerator application.RoomNumberGenerator,
	encryptor application.Encryptor,
	roomSaver application.RoomSaver,
	eventSaver event.UnpublishedSaver,
	unitOfWork kernel.UnitOfWork,
) CreateRoomUseCase {
	return &createRoomUseCase{
		eventIDGenerator: eventIDGenerator,
		roomIDGenerator:  roomIDGenerator,
		numberGenerator:  numberGenerator,
		encryptor:        encryptor,
		roomSaver:        roomSaver,
		eventSaver:       eventSaver,
		unitOfWork:       unitOfWork,
	}
}

func (uc *createRoomUseCase) Execute(ctx context.Context, input *CreateRoomInput) (*CreateRoomOutput, error) {
	now := time.Now().UTC()
	var passwordEncrypted string
	if len(input.Password) != 0 {
		encrypted, err := uc.encryptor.Encrypt(input.Password.String())
		if err != nil {
			return nil, err
		}
		passwordEncrypted = encrypted
	}

	room := domain.NewRoom(
		uc.roomIDGenerator.Generate(),
		input.Owner,
		uc.numberGenerator.Generate(),
		passwordEncrypted,
		make([]kernel.UserID, 0, 1),
		input.MaxMemberCount,
		now,
	)

	if err := room.Create(uc.eventIDGenerator); err != nil {
		return nil, err
	}

	if err := room.Join(input.Owner, uc.eventIDGenerator); err != nil {
		return nil, err
	}

	if err := uc.unitOfWork.Execute(ctx, func(txCtx context.Context) error {
		if err := uc.roomSaver.Save(txCtx, room); err != nil {
			return err
		}

		if err := uc.eventSaver.Saves(txCtx, room.GetEvents()); err != nil {
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
