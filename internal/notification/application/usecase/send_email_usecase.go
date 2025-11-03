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
	Theme          domain.NoticeTheme
	Title          string
	Content        string
}

func (i *SendEmailInput) Validate() error {
	if err := i.RecipientEmail.Validate(); err != nil {
		return err
	}
	return nil
}

type notifyUseCase struct {
	emailNotifier application.EmailNotifier
	saver         application.NoticeSaver
	IDGenerator   application.NoticeIDGenerator
}

func NewSendEmailUseCase(
	emailNotifier application.EmailNotifier,
	saver application.NoticeSaver,
	IDGenerator application.NoticeIDGenerator,
) SendEmailUseCase {
	return &notifyUseCase{emailNotifier: emailNotifier, saver: saver, IDGenerator: IDGenerator}
}

func (uc *notifyUseCase) Execute(ctx context.Context, input *SendEmailInput) (*kernel.NoOutput, error) {
	notice := domain.NewNotice(uc.IDGenerator.Generate(), input.Recipient, input.Theme, input.Title, input.Content, domain.Contact(input.RecipientEmail))

	if err := uc.emailNotifier.Enqueue(ctx, input.RecipientEmail, notice, func() error {
		if err := uc.saver.Save(context.Background(), notice); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return nil, nil
}
