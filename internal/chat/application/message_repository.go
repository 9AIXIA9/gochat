package application

import (
	"context"
	"gochat/internal/chat/domain"
)

type MessagesSaver interface {
	Save(ctx context.Context, messages []*domain.Message) error
}
