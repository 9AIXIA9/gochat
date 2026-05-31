package command

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/command"
	"gochat/internal/shared/kernel"
)

func NewSendRoomMessageCommandHandler(
	uc application.SendRoomMessageUseCase,
) command.Handler {
	return command.AdaptUsecaseToCommandHandler(
		uc,
		domain.ToSendRoomMessageCommand,
		func(com *domain.SendRoomMessageCommand) *application.SendRoomMessageInput {
			return &application.SendRoomMessageInput{
				SenderID: kernel.UserID(com.AggregateID()),
				RoomID:   com.RoomID(),
				Content:  com.Content(),
			}
		},
		nil,
		nil,
	)
}
