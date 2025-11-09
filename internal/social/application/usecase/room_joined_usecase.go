package usecase

import (
	"context"
	chatDomain "gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	socialDomain "gochat/internal/social/domain"
)

type RoomJoinedUseCase kernel.UseCase[*RoomJoinedInput, *kernel.NoOutput]

type RoomJoinedInput struct {
	RoomID socialDomain.RoomID
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
	publisher   event.Publisher
}

func NewRoomJoinedUseCase(
	idGenerator event.IDGenerator,
	publisher event.Publisher,
) RoomJoinedUseCase {
	return &roomJoinedUseCase{
		idGenerator: idGenerator,
		publisher:   publisher,
	}
}

func (uc *roomJoinedUseCase) Execute(_ context.Context, input *RoomJoinedInput) (*kernel.NoOutput, error) {
	notificationEv, err := notificationDomain.NewRoomJoinedEvent(uc.idGenerator.Generate(), notificationDomain.RoomID(input.RoomID), input.UserID)
	if err != nil {
		return nil, err
	}
	if err := uc.publisher.Publish(notificationEv); err != nil {
		return nil, err
	}

	chatEv, err := chatDomain.NewRoomJoinedEvent(uc.idGenerator.Generate(), chatDomain.RoomID(input.RoomID), input.UserID)
	if err != nil {
		return nil, err
	}
	if err := uc.publisher.Publish(chatEv); err != nil {
		return nil, err
	}
	return nil, nil
}
