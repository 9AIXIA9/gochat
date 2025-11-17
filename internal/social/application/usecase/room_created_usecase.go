package usecase

import (
	"context"
	chatDomain "gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/application"
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
	publisher   event.SinglePublisher
	finder      application.RoomFinderByID
}

func NewRoomCreatedUseCase(
	idGenerator event.IDGenerator,
	publisher event.SinglePublisher,
	finder application.RoomFinderByID,
) RoomCreatedUseCase {
	return &roomCreatedUseCase{
		idGenerator: idGenerator,
		publisher:   publisher,
		finder:      finder,
	}
}

func (uc *roomCreatedUseCase) Execute(ctx context.Context, input *RoomCreatedInput) (*kernel.NoOutput, error) {
	room, err := uc.finder.FindByID(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	chatEv, err := chatDomain.NewRoomCreatedEvent(uc.idGenerator.Generate(), room.ID(), room.Number())
	if err != nil {
		return nil, err
	}
	if err := uc.publisher.Publish(chatEv); err != nil {
		return nil, err
	}
	return nil, nil
}
