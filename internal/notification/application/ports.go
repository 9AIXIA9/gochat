package application

import (
	"context"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
)

type EmailNotifier interface {
	Enqueue(ctx context.Context, email kernel.Email, mail *domain.Mail, onSuccess func() error) error
}

type MessageNotifier interface {
	Notify(ctx context.Context, message *domain.Message) error
}

type MailIDGenerator interface {
	Generate() domain.MailID
}
