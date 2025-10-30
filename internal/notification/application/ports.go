package application

import (
	"context"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
)

type EmailNotifier interface {
	Enqueue(ctx context.Context, email kernel.Email, notice *domain.Notice, onSuccess func() error) error
}

type NoticeIDGenerator interface {
	Generate() domain.NoticeID
}
