package usecase

import (
	"context"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/application"
	"gochat/internal/social/domain"
)

type JoinRoomUseCase kernel.UseCase[*JoinRoomInput, *kernel.NoOutput]

type JoinRoomInput struct {
	UserID     kernel.UserID
	RoomNumber domain.RoomNumber
	Password   domain.Password
}

func (r *JoinRoomInput) Validate() error {
	if err := r.Password.Validate(); err != nil {
		return err
	}
	if len(r.UserID) == 0 || len(r.RoomNumber) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type joinRoomUseCase struct {
	eventIDGenerator event.IDGenerator
	finder           application.RoomFinderByNumber
	comparator       application.Comparator
	roomMemberSaver  application.RoomMemberSaver
	eventSaver       event.UnpublishedSaver
}

func NewJoinRoomUseCase(
	eventIDGenerator event.IDGenerator,
	eventSaver event.UnpublishedSaver,
	finder application.RoomFinderByNumber,
	comparator application.Comparator,
	roomMemberSaver application.RoomMemberSaver,
) JoinRoomUseCase {
	return &joinRoomUseCase{
		eventIDGenerator: eventIDGenerator,
		finder:           finder,
		comparator:       comparator,
		roomMemberSaver:  roomMemberSaver,
		eventSaver:       eventSaver,
	}
}

func (uc *joinRoomUseCase) Execute(ctx context.Context, input *JoinRoomInput) (*kernel.NoOutput, error) {
	room, err := uc.finder.FindByNumber(ctx, input.RoomNumber)
	if err != nil {
		return nil, err
	}

	if encrypted := room.PasswordEncrypted(); len(encrypted) != 0 {
		if err := uc.comparator.Compare(encrypted, input.Password.String()); err != nil {
			return nil, err
		}
	}

	if err := room.Join(input.UserID, uc.eventIDGenerator); err != nil {
		return nil, err
	}

	if err := uc.roomMemberSaver.SaveMember(ctx, room.ID(), input.UserID); err != nil {
		return nil, err
	}

	if err := uc.eventSaver.Saves(ctx, room.GetEvents()); err != nil {
		return nil, err
	}

	return nil, nil
}
