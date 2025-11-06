package usecase

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

var _ SendRoomMessageUseCase = (*sendRoomMessageUseCase)(nil)

type SendRoomMessageUseCase kernel.UseCase[*SendRoomMessageInput, *kernel.NoOutput]

type SendRoomMessageInput struct {
	SenderID kernel.UserID
	RoomID   domain.RoomID //TODO 通过Number来指定接收者 而不是ID
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
	messageSaver          application.MessageSaver
	unpublishedEventSaver event.UnpublishedSaver
}

func NewSendRoomMessageUseCase(
	messageIDGenerator application.MessageIDGenerator,
	eventIDGenerator event.IDGenerator,
	roomMembersFinder application.RoomMembersFinder,
	messageSaver application.MessageSaver,
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
	members, err := uc.roomMembersFinder.FindMembersByRoomID(ctx, input.RoomID)
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

	room := domain.NewRoom(input.RoomID, members, make([]*domain.Message, 0, 1))

	if err := room.SendMessage(uc.messageIDGenerator.Generate(), input.SenderID, input.Content, uc.eventIDGenerator); err != nil {
		return nil, err
	}

	if err := uc.messageSaver.Saves(ctx, room.Messages()); err != nil {
		return nil, err
	}

	events := room.GetEvents()
	if err := uc.unpublishedEventSaver.Saves(ctx, events); err != nil {
		return nil, err
	}
	return nil, nil
}
