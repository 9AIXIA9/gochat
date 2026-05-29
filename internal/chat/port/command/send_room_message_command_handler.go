package command

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewSendRoomMessageCommandHandler(uc application.SendRoomMessageUseCase) event.Handler {
	return event.AdaptUsecaseToHandler(
		uc,
		domain.ToSendRoomMessageCommand,
		func(ev *domain.SendRoomMessageCommand) *application.SendRoomMessageInput {
			return &application.SendRoomMessageInput{
				SenderID: kernel.UserID(ev.AggregateID()),
				RoomID:   ev.RoomID(),
				Content:  ev.Content(),
			}
		},
		nil,
		nil,
	)
}
