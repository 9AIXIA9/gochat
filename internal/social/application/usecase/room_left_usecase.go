package usecase

import (
	"context"
	chatDomain "gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type RoomLeftUseCase kernel.UseCase[*RoomLeftInput, *kernel.NoOutput]

type RoomLeftInput struct {
	RoomID kernel.RoomID
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
	saver       event.UnpublishedEventSaver
}

func NewRoomLeftUseCase(
	idGenerator event.IDGenerator,
	saver event.UnpublishedEventSaver,
) RoomLeftUseCase {
	return &roomLeftUseCase{
		idGenerator: idGenerator,
		saver:       saver,
	}
}

func (uc *roomLeftUseCase) Execute(ctx context.Context, input *RoomLeftInput) (*kernel.NoOutput, error) {
	chatEv, err := chatDomain.NewRoomLeftEvent(uc.idGenerator.Generate(), input.RoomID, input.UserID)
	if err != nil {
		return nil, err
	}
	if err := uc.saver.Save(ctx, chatEv); err != nil {
		return nil, err
	}
	return nil, nil
}
