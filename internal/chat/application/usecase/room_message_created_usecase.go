package usecase

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

var _ RoomMessageCreatedUseCase = (*roomMessageCreatedUseCase)(nil)

type RoomMessageCreatedUseCase kernel.UseCase[*RoomMessageCreatedInput, *kernel.NoOutput]

type RoomMessageCreatedInput struct {
	RoomID    domain.RoomID
	MessageID domain.MessageID
}

func (i *RoomMessageCreatedInput) Validate() error {
	if len(i.RoomID) == 0 || len(i.MessageID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type roomMessageCreatedUseCase struct {
	eventIDGenerator  event.IDGenerator
	publisher         event.Publisher
	roomMessageFinder application.RoomMessageFinder
	roomMemberFinder  application.RoomMemberFinder
}

func NewRoomMessageCreatedUseCase(
	eventIDGenerator event.IDGenerator,
	roomMessageFinder application.RoomMessageFinder,
	roomMemberFinder application.RoomMemberFinder,
	publisher event.Publisher,
) RoomMessageCreatedUseCase {
	return &roomMessageCreatedUseCase{
		eventIDGenerator:  eventIDGenerator,
		publisher:         publisher,
		roomMessageFinder: roomMessageFinder,
		roomMemberFinder:  roomMemberFinder,
	}
}

func (uc *roomMessageCreatedUseCase) Execute(ctx context.Context, input *RoomMessageCreatedInput) (*kernel.NoOutput, error) {
	message, err := uc.roomMessageFinder.FindRoomMessage(ctx, input.RoomID, input.MessageID)
	if err != nil {
		return nil, err
	}

	members, err := uc.roomMemberFinder.FindMember(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	ev, err := notificationDomain.NewRoomMessageCreatedEvent(
		uc.eventIDGenerator.Generate(),
		notificationDomain.RoomID(input.RoomID),
		notificationDomain.MessageID(message.ID()),
		members,
		message.Sender(),
		message.Content(),
		message.SentAt(),
	)
	if err != nil {
		return nil, err
	}

	if err := uc.publisher.Publish(ev); err != nil {
		return nil, err
	}

	return nil, nil
}
