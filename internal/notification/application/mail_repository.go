package application

import (
	"context"
	"gochat/internal/notification/domain"
)

type MailRepository interface {
	MailSaver
}

type MailSaver interface {
	Save(ctx context.Context, mail *domain.Mail) error
}
