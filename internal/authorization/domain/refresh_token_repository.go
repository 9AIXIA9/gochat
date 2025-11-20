package domain

import (
	"context"
)

type RefreshTokenRepository interface {
	RefreshTokenSaver
	RefreshTokenFinder
}

type RefreshTokenSaver interface {
	Save(ctx context.Context, token *RefreshTokenEntity) error
}

type RefreshTokenFinder interface {
	FindByToken(ctx context.Context, token RefreshToken) (*RefreshTokenEntity, error)
}
