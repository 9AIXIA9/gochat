package application

import (
	"context"
	"gochat/internal/notification/domain"
)

type MailSaver interface {
	Save(ctx context.Context, mail *domain.Mail) error
}
