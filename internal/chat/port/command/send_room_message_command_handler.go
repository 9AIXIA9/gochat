package command

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/contract"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewSendRoomMessageCommandHandler(
	uc application.SendRoomMessageUseCase,
	gateway contract.GatewayService,
) event.Handler {
	return event.AdaptUsecaseToCommandHandler(
		uc,
		domain.ToSendRoomMessageCommand,
		func(ev *domain.SendRoomMessageCommand) *application.SendRoomMessageInput {
			return &application.SendRoomMessageInput{
				SenderID: kernel.UserID(ev.AggregateID()),
				RoomID:   ev.RoomID(),
				Content:  ev.Content(),
			}
		},
		func(ctx context.Context, action *domain.SendRoomMessageCommand, output *kernel.NoOutput) {
			_ = gateway.PushToUser(ctx, kernel.UserID(action.AggregateID()), []byte(`{"type":"send_room_message_succeeded"}`))
		},
		func(ctx context.Context, action *domain.SendRoomMessageCommand, err error) {
			_ = gateway.PushToUser(ctx, kernel.UserID(action.AggregateID()), []byte(`{"type":"send_room_message_failed"}`))
		},
	)
}
