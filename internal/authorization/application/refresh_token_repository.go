package application

import (
	"context"
	"gochat/internal/authorization/domain"
)

type RefreshTokenSaver interface {
	Save(ctx context.Context, token *domain.RefreshTokenEntity) error
}

type RefreshTokenFinder interface {
	FindByToken(ctx context.Context, token domain.RefreshToken) (*domain.RefreshTokenEntity, error)
}
