package kafka

import (
	"context"
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"
)

type UserFinderByID interface {
	FindByID(ctx context.Context, id kernel.UserID) (*domain.User, error)
}
