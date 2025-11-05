package usecase

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

var _ SendRoomMessageUseCase = (*sendRoomMessageUseCase)(nil)

type SendRoomMessageUseCase kernel.UseCase[*SendRoomMessageInput, *kernel.NoOutput]

type SendRoomMessageInput struct {
	SenderID kernel.UserID
	RoomID   domain.RoomID
	Content  string
}

func (i *SendRoomMessageInput) Validate() error {
	if len(i.Content) == 0 {
		return myErrors.ErrEmptyInput
	}
	if len(i.SenderID) == 0 || len(i.RoomID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type sendRoomMessageUseCase struct {
	messageIDGenerator    application.MessageIDGenerator
	eventIDGenerator      event.IDGenerator
	roomMembersFinder     application.RoomMembersFinder
	messageSaver          application.RoomMessagesSaver
	unpublishedEventSaver event.UnpublishedSaver
}

func NewSendRoomMessageUseCase(
	messageIDGenerator application.MessageIDGenerator,
	eventIDGenerator event.IDGenerator,
	roomMembersFinder application.RoomMembersFinder,
	messageSaver application.RoomMessagesSaver,
	unpublishedEventSaver event.UnpublishedSaver,
) SendRoomMessageUseCase {
	return &sendRoomMessageUseCase{
		messageIDGenerator:    messageIDGenerator,
		eventIDGenerator:      eventIDGenerator,
		roomMembersFinder:     roomMembersFinder,
		messageSaver:          messageSaver,
		unpublishedEventSaver: unpublishedEventSaver,
	}
}

func (uc *sendRoomMessageUseCase) Execute(ctx context.Context, input *SendRoomMessageInput) (*kernel.NoOutput, error) {
	members, err := uc.roomMembersFinder.FindByRoomID(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	states := make([]*domain.RecipientMessageState, 0, len(members))
	for _, member := range members {
		state := domain.NewRecipientMessageState(
			member,
			domain.MessageStateCreated,
		)
		states = append(states, state)
	}

	message := domain.NewRoomMessage(
		uc.messageIDGenerator.Generate(),
		input.RoomID,
		input.SenderID,
		input.Content,
		time.Now().UTC(),
		states,
	)

	user := domain.NewUser(input.SenderID, nil, make([]*domain.RoomMessage, 0, 1))

	if err := user.SendRoomMessage(message, uc.eventIDGenerator); err != nil {
		return nil, err
	}

	if err := uc.messageSaver.SaveRoomMessages(ctx, user.RoomMessages()); err != nil {
		return nil, err
	}

	events := user.GetEvents()
	if err := uc.unpublishedEventSaver.Saves(ctx, events); err != nil {
		return nil, err
	}
	return nil, nil
}
