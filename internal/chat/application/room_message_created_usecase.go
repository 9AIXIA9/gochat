package application

import (
	"context"
	"errors"
	"gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
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
	roomFinder        domain.RoomshipsFinderByRoomID
}

func NewRoomMessageCreatedUseCase(
	eventIDGenerator event.IDGenerator,
	roomMessageFinder domain.RoomMessageFinder,
	roomFinder domain.RoomshipsFinderByRoomID,
	creator event.UnpublishedEventCreator,
) (RoomMessageCreatedUseCase, error) {
	if err := utils.CheckInterfaces(
		eventIDGenerator,
		roomMessageFinder,
		roomFinder,
		creator,
	); err != nil {
		return nil, err
	}
	return &roomMessageCreatedUseCase{
		eventIDGenerator:  eventIDGenerator,
		creator:           creator,
		roomMessageFinder: roomMessageFinder,
		roomFinder:        roomFinder,
	}, nil
}

func (uc *roomMessageCreatedUseCase) Execute(ctx context.Context, input *RoomMessageCreatedInput) (*kernel.NoOutput, error) {
	message, err := uc.roomMessageFinder.FindRoomMessage(ctx, input.MessageID)
	if err != nil {
		return nil, err
	}

	roomships, err := uc.roomFinder.FindsByRoomID(ctx, message.RoomID())
	if err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return nil, nil

		}
		return nil, err
	}

	recipients := make([]kernel.UserID, 0, len(roomships))
	for _, roomship := range roomships {
		if roomship.UserID() != message.SenderID() {
			recipients = append(recipients, roomship.UserID())
		}
	}

	if len(recipients) == 0 {
		return nil, nil
	}

	ev, err := notificationDomain.NewRoomMessageNotificationRequestedEvent(
		message.ID(),
		message.SenderID(),
		message.RoomID(),
		recipients,
		message.Content(),
		message.SentAt(),
		uc.eventIDGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.CreateUnpublishedEvent(ctx, ev); err != nil {
		return nil, err
	}

	return nil, nil
}
