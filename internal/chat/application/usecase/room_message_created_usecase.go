package usecase

import (
	"context"
	"gochat/internal/chat/application"
	notificationDomain "gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

var _ RoomMessageCreatedUseCase = (*roomMessageCreatedUseCase)(nil)

type RoomMessageCreatedUseCase kernel.UseCase[*RoomMessageCreatedInput, *kernel.NoOutput]

type RoomMessageCreatedInput struct {
	RoomID     kernel.RoomID
	MessageID  kernel.MessageID
	Recipients []kernel.UserID
}

func (i *RoomMessageCreatedInput) Validate() error {
	if len(i.RoomID) == 0 || len(i.MessageID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type roomMessageCreatedUseCase struct {
	eventIDGenerator  event.IDGenerator
	publisher         event.ManyPublisher
	roomMessageFinder application.RoomMessageFinder
}

func NewRoomMessageCreatedUseCase(
	eventIDGenerator event.IDGenerator,
	roomMessageFinder application.RoomMessageFinder,
	publisher event.ManyPublisher,
) RoomMessageCreatedUseCase {
	return &roomMessageCreatedUseCase{
		eventIDGenerator:  eventIDGenerator,
		publisher:         publisher,
		roomMessageFinder: roomMessageFinder,
	}
}

func (uc *roomMessageCreatedUseCase) Execute(ctx context.Context, input *RoomMessageCreatedInput) (*kernel.NoOutput, error) {
	if len(input.Recipients) == 0 {
		//没有接收者不发送通知
		return nil, nil
	}

	message, err := uc.roomMessageFinder.FindRoomMessage(ctx, input.RoomID, input.MessageID)
	if err != nil {
		return nil, err
	}

	evs := make([]event.Event, 0, len(input.Recipients))
	for _, recipient := range input.Recipients {
		ev, err := notificationDomain.NewMessageNotificationRequestedEvent(
			uc.eventIDGenerator.Generate(),
			message.ID(),
			recipient,
			message.Sender(),
			message.Content(),
			message.SentAt(),
		)
		if err != nil {
			return nil, err
		}
		evs = append(evs, ev)
	}

	if err := uc.publisher.Publishes(evs); err != nil {
		return nil, err
	}

	return nil, nil
}
