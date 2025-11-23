package domain

import (
	"context"
)

type RefreshTokenRepository interface {
	RefreshTokenCreator
	RefreshTokenFinder
}

type RefreshTokenCreator interface {
	Create(ctx context.Context, token *RefreshTokenEntity) error
}

type RefreshTokenFinder interface {
	FindByToken(ctx context.Context, token RefreshToken) (*RefreshTokenEntity, error)
}
