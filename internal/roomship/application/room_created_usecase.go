package application

import (
	"context"
	chatDomain "gochat/internal/chat/domain"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
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
	idGenerator event.IDGenerator
	creator     event.UnpublishedEventCreator
	finder      domain.RoomFinderByID
}

func NewRoomCreatedUseCase(
	idGenerator event.IDGenerator,
	creator event.UnpublishedEventCreator,
	finder domain.RoomFinderByID,
) RoomCreatedUseCase {
	return &roomCreatedUseCase{
		idGenerator: idGenerator,
		creator:     creator,
		finder:      finder,
	}
}

func (uc *roomCreatedUseCase) Execute(ctx context.Context, input *RoomCreatedInput) (*kernel.NoOutput, error) {
	room, err := uc.finder.FindByID(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	chatEv, err := chatDomain.NewRoomCreatedEvent(room.ID(), room.OwnerID(), uc.idGenerator)
	if err != nil {
		return nil, err
	}
	if err := uc.creator.CreateUnpublishedEvent(ctx, chatEv); err != nil {
		return nil, err
	}
	return nil, nil
}
