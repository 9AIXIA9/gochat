package command

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"

	"gochat/internal/shared/command"
)

func NewSendPrivateMessageCommandHandler(
	uc application.SendPrivateMessageUseCase,
	publisher command.ReceiptAsyncPublisher,
	generator command.ReceiptIDGenerator,
) command.Handler {
	return command.AdaptUsecaseToCommandHandler(
		uc,
		domain.ToSendPrivateMessageCommand,
		func(com *domain.SendPrivateMessageCommand) *application.SendPrivateMessageInput {
			return &application.SendPrivateMessageInput{
				SenderID:    kernel.UserID(com.AggregateID()),
				RecipientID: com.RecipientID(),
				Content:     com.Content(),
			}
		},
		publisher,
		generator,
	)
}
