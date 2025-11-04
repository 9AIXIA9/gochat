package usecase

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
)

type SendEmailUseCase kernel.UseCase[*SendEmailInput, *kernel.NoOutput]

type SendEmailInput struct {
	Recipient      kernel.UserID
	RecipientEmail kernel.Email
	Theme          domain.MailTheme
	Title          string
	Content        string
}

func (i *SendEmailInput) Validate() error {
	if err := i.RecipientEmail.Validate(); err != nil {
		return err
	}
	return nil
}

type sendEmailUseCase struct {
	emailNotifier application.EmailNotifier
	saver         application.MailSaver
	IDGenerator   application.MailIDGenerator
}

func NewSendEmailUseCase(
	emailNotifier application.EmailNotifier,
	saver application.MailSaver,
	IDGenerator application.MailIDGenerator,
) SendEmailUseCase {
	return &sendEmailUseCase{emailNotifier: emailNotifier, saver: saver, IDGenerator: IDGenerator}
}

func (uc *sendEmailUseCase) Execute(ctx context.Context, input *SendEmailInput) (*kernel.NoOutput, error) {
	mail := domain.NewMail(uc.IDGenerator.Generate(), input.Recipient, input.Theme, input.Title, input.Content, input.RecipientEmail)

	if err := uc.emailNotifier.Enqueue(ctx, input.RecipientEmail, mail, func() error {
		if err := uc.saver.Save(context.Background(), mail); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return nil, nil
}
