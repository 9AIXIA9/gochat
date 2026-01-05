package application

import (
	"context"
	"errors"
	chatDomain "gochat/internal/chat/domain"
	profileDomain "gochat/internal/profile/domain"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type RoomCreatedUseCase kernel.UseCase[*RoomCreatedInput, *kernel.NoOutput]

type RoomCreatedInput struct {
	RoomID kernel.RoomID
}

func (r *RoomCreatedInput) Validate() error {
	if len(r.RoomID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "room id is empty")
	}

	return nil
}

type roomCreatedUseCase struct {
	roomshipIDGenerator domain.RoomshipIDGenerator
	idGenerator         event.IDGenerator
	finder              domain.RoomFinderByID
	roomshipCreator     domain.RoomshipCreator
	eventCreator        event.UnpublishedEventsCreator
}

func NewRoomCreatedUseCase(
	roomshipIDGenerator domain.RoomshipIDGenerator,
	idGenerator event.IDGenerator,
	finder domain.RoomFinderByID,
	roomshipCreator domain.RoomshipCreator,
	eventCreator event.UnpublishedEventsCreator,
) (RoomCreatedUseCase, error) {
	if err := validate.NotNil(
		roomshipIDGenerator,
		idGenerator,
		finder,
		roomshipCreator,
		eventCreator,
	); err != nil {
		return nil, err
	}

	return &roomCreatedUseCase{
		roomshipIDGenerator: roomshipIDGenerator,
		idGenerator:         idGenerator,
		finder:              finder,
		roomshipCreator:     roomshipCreator,
		eventCreator:        eventCreator,
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

	if err := uc.roomshipCreator.Create(ctx, roomship); err != nil {
		if errors.Is(err, myErrors.ErrDuplicatedKey) {
			return nil, nil
		}
		return nil, err
	}

	profileEvCreated, err := profileDomain.NewRoomCreatedEvent(room.ID(), room.CreatedAt(), uc.idGenerator)
	if err != nil {
		return nil, err
	}

	chatEvCreated, err := chatDomain.NewRoomCreatedEvent(room.ID(), uc.idGenerator)
	if err != nil {
		return nil, err
	}

	if err := uc.eventCreator.CreateUnpublishedEvents(ctx, []event.Event{
		profileEvCreated,
		chatEvCreated,
	}); err != nil {
		return nil, err
	}

	return nil, nil
}
