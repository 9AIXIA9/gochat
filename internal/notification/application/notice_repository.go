package application

import (
	"context"
	"gochat/internal/notification/domain"
)

type NoticeRepository interface {
	NoticeSaver
}

type NoticeSaver interface {
	Save(ctx context.Context, notice *domain.Notice) error
}
