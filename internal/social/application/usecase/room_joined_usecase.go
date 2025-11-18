package usecase

import (
	"context"
	chatDomain "gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type RoomJoinedUseCase kernel.UseCase[*RoomJoinedInput, *kernel.NoOutput]

type RoomJoinedInput struct {
	RoomID kernel.RoomID
	UserID kernel.UserID
}

func (r *RoomJoinedInput) Validate() error {
	if len(r.RoomID) == 0 || len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type roomJoinedUseCase struct {
	idGenerator event.IDGenerator
	saver       event.UnpublishedEventSaver
}

func NewRoomJoinedUseCase(
	idGenerator event.IDGenerator,
	saver event.UnpublishedEventSaver,
) RoomJoinedUseCase {
	return &roomJoinedUseCase{
		idGenerator: idGenerator,
		saver:       saver,
	}
}

func (uc *roomJoinedUseCase) Execute(ctx context.Context, input *RoomJoinedInput) (*kernel.NoOutput, error) {
	chatEv, err := chatDomain.NewRoomJoinedEvent(uc.idGenerator.Generate(), input.RoomID, input.UserID)
	if err != nil {
		return nil, err
	}
	if err := uc.saver.Save(ctx, chatEv); err != nil {
		return nil, err
	}
	return nil, nil
}
