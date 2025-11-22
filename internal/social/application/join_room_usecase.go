package application

import (
	"context"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"
)

type JoinRoomUseCase kernel.UseCase[*JoinRoomInput, *kernel.NoOutput]

type JoinRoomInput struct {
	UserID     kernel.UserID
	RoomNumber kernel.RoomNumber
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
	finder           domain.RoomFinderByNumber
	comparator       domain.Comparator
	roomMemberSaver  domain.RoomMemberSaver
	eventsCreator    event.UnpublishedEventsCreator
	unitOfWork       kernel.UnitOfWork
}

func NewJoinRoomUseCase(
	eventIDGenerator event.IDGenerator,
	eventsCreator event.UnpublishedEventsCreator,
	finder domain.RoomFinderByNumber,
	comparator domain.Comparator,
	roomMemberSaver domain.RoomMemberSaver,
	unitOfWork kernel.UnitOfWork,
) JoinRoomUseCase {
	return &joinRoomUseCase{
		eventIDGenerator: eventIDGenerator,
		finder:           finder,
		comparator:       comparator,
		roomMemberSaver:  roomMemberSaver,
		eventsCreator:    eventsCreator,
		unitOfWork:       unitOfWork,
	}
}

func (uc *joinRoomUseCase) Execute(ctx context.Context, input *JoinRoomInput) (*kernel.NoOutput, error) {
	room, err := uc.finder.FindByNumber(ctx, input.RoomNumber)
	if err != nil {
		return nil, err
	}

	if err := room.AddMember(
		input.UserID,
		input.Password,
		uc.comparator,
		uc.eventIDGenerator,
	); err != nil {
		return nil, err
	}

	if err := uc.unitOfWork.Execute(ctx, func(txCtx context.Context) error {
		if err := uc.roomMemberSaver.SaveMember(txCtx, room.ID(), input.UserID); err != nil {
			return err
		}

		if err := uc.eventsCreator.CreateUnpublishedEvents(txCtx, room.GetEvents()); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}
