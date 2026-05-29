package command

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"

	"gochat/internal/shared/event"
)

func NewSendPrivateMessageCommandHandler(uc application.SendPrivateMessageUseCase) event.Handler {
	return event.AdaptUsecaseToHandler(
		uc,
		domain.ToSendPrivateMessageCommand,
		func(ev *domain.SendPrivateMessageCommand) *application.SendPrivateMessageInput {
			return &application.SendPrivateMessageInput{
				SenderID:    kernel.UserID(ev.AggregateID()),
				RecipientID: ev.RecipientID(),
				Content:     ev.Content(),
			}
		},
		nil,
		nil,
	)
}
