package application

import (
	"context"
	chatDomain "gochat/internal/chat/domain"
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
	if err := utils.CheckInterfaces(
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
	//TODO 可能重复的原因在于 事件被重复发布了，需要处理幂等性 所以要为仓库添加唯一性索引
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
		return nil, err
	}

	chatEvCreated, err := chatDomain.NewRoomCreatedEvent(room.ID(), uc.idGenerator)
	if err != nil {
		return nil, err
	}

	if err := uc.eventCreator.CreateUnpublishedEvents(ctx, []event.Event{chatEvCreated}); err != nil {
		return nil, err
	}

	return nil, nil
}
