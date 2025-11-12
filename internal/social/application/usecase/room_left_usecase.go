package usecase

import (
	"context"
	chatDomain "gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	socialDomain "gochat/internal/social/domain"
)

type RoomLeftUseCase kernel.UseCase[*RoomLeftInput, *kernel.NoOutput]

type RoomLeftInput struct {
	RoomID socialDomain.RoomID
	UserID kernel.UserID
}

func (r *RoomLeftInput) Validate() error {
	if len(r.RoomID) == 0 || len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type roomLeftUseCase struct {
	idGenerator event.IDGenerator
	publisher   event.Publisher
}

func NewRoomLeftUseCase(
	idGenerator event.IDGenerator,
	publisher event.Publisher,
) RoomLeftUseCase {
	return &roomLeftUseCase{
		idGenerator: idGenerator,
		publisher:   publisher,
	}
}

func (uc *roomLeftUseCase) Execute(_ context.Context, input *RoomLeftInput) (*kernel.NoOutput, error) {
	chatEv, err := chatDomain.NewRoomLeftEvent(uc.idGenerator.Generate(), chatDomain.RoomID(input.RoomID), input.UserID)
	if err != nil {
		return nil, err
	}
	if err := uc.publisher.Publish(chatEv); err != nil {
		return nil, err
	}
	return nil, nil
}
