//go:generate mockgen -source=refresh_token_repository.go -destination=./mocks/mock_refresh_token_repository.go -package=mocks
package domain

import (
	"context"
)

type RefreshTokenRepository interface {
	RefreshTokenUpserter
	RefreshTokenFinder
}

type RefreshTokenUpserter interface {
	Upsert(ctx context.Context, token *RefreshTokenEntity) error
}

type RefreshTokenFinder interface {
	FindByToken(ctx context.Context, token RefreshToken) (*RefreshTokenEntity, error)
}
