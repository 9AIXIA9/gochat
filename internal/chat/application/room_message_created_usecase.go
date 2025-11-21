package application

import (
	"context"
	"gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

var _ RoomMessageCreatedUseCase = (*roomMessageCreatedUseCase)(nil)

type RoomMessageCreatedUseCase kernel.UseCase[*RoomMessageCreatedInput, *kernel.NoOutput]

type RoomMessageCreatedInput struct {
	MessageID kernel.MessageID
}

func (i *RoomMessageCreatedInput) Validate() error {
	if len(i.MessageID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type roomMessageCreatedUseCase struct {
	eventIDGenerator  event.IDGenerator
	saver             event.UnpublishedEventsSaver
	roomMessageFinder domain.RoomMessageFinder
	roomMembersFinder domain.RoomMembersFinder
}

func NewRoomMessageCreatedUseCase(
	eventIDGenerator event.IDGenerator,
	roomMessageFinder domain.RoomMessageFinder,
	roomMembersFinder domain.RoomMembersFinder,
	saver event.UnpublishedEventsSaver,
) RoomMessageCreatedUseCase {
	return &roomMessageCreatedUseCase{
		eventIDGenerator:  eventIDGenerator,
		saver:             saver,
		roomMessageFinder: roomMessageFinder,
		roomMembersFinder: roomMembersFinder,
	}
}

func (uc *roomMessageCreatedUseCase) Execute(ctx context.Context, input *RoomMessageCreatedInput) (*kernel.NoOutput, error) {
	message, err := uc.roomMessageFinder.FindRoomMessage(ctx, input.MessageID)
	if err != nil {
		return nil, err
	}

	members, err := uc.roomMembersFinder.FindMembers(ctx, message.RoomID())
	if err != nil {
		return nil, err
	}

	evs := make([]event.Event, 0, len(members))
	for _, member := range members {
		if member == message.SenderID() {
			continue
		}
		ev, err := notificationDomain.NewMessageNotificationRequestedEvent(
			uc.eventIDGenerator.Generate(),
			message.ID(),
			member,
			message.SenderID(),
			message.Content(),
			message.SentAt(),
		)
		if err != nil {
			return nil, err
		}
		evs = append(evs, ev)
	}

	if err := uc.saver.Saves(ctx, evs); err != nil {
		return nil, err
	}

	return nil, nil
}
