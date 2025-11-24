package application

import (
	"context"
	"errors"
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
	creator           event.UnpublishedEventCreator
	roomMessageFinder domain.RoomMessageFinder
	roomFinder        domain.RoomFinderByID
}

func NewRoomMessageCreatedUseCase(
	eventIDGenerator event.IDGenerator,
	roomMessageFinder domain.RoomMessageFinder,
	roomFinder domain.RoomFinderByID,
	creator event.UnpublishedEventCreator,
) RoomMessageCreatedUseCase {
	return &roomMessageCreatedUseCase{
		eventIDGenerator:  eventIDGenerator,
		creator:           creator,
		roomMessageFinder: roomMessageFinder,
		roomFinder:        roomFinder,
	}
}

func (uc *roomMessageCreatedUseCase) Execute(ctx context.Context, input *RoomMessageCreatedInput) (*kernel.NoOutput, error) {
	message, err := uc.roomMessageFinder.FindRoomMessage(ctx, input.MessageID)
	if err != nil {
		return nil, err
	}

	room, err := uc.roomFinder.FindByID(ctx, message.RoomID())
	if err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return nil, nil

		}
		return nil, err
	}

	recipients := room.Members()
	for i, recipient := range recipients {
		if recipient == message.SenderID() {
			recipients = append(recipients[:i], recipients[i+1:]...)
			break
		}
	}

	if len(recipients) == 0 {
		return nil, nil
	}

	ev, err := notificationDomain.NewRoomMessageNotificationRequestedEvent(
		uc.eventIDGenerator.Generate(),
		message.ID(),
		message.SenderID(),
		message.RoomID(),
		recipients,
		message.Content(),
		message.SentAt(),
	)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.CreateUnpublishedEvent(ctx, ev); err != nil {
		return nil, err
	}

	return nil, nil
}
