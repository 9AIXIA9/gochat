package kafka

import (
	"context"
	"fmt"
	"gochat/internal/notification/application/usecase"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewEmailSendingRequestedHandler(uc usecase.SendEmailUseCase) event.Handler {
	return func(ctx context.Context, e event.Event) error {
		fmt.Println("email sending requested event received")

		emailSendingRequested, err := domain.ToEmailSendingRequestedEvent(e)
		if err != nil {
			return err
		}

		_, err = uc.Execute(ctx, &usecase.SendEmailInput{
			Recipient:      kernel.UserID(emailSendingRequested.AggregateID()),
			RecipientEmail: emailSendingRequested.Email(),
			Theme:          emailSendingRequested.Theme(),
			Title:          emailSendingRequested.Title(),
			Content:        emailSendingRequested.Content(),
		})
		return err
	}
}
