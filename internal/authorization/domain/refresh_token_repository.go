package domain

import (
	"context"
)

type RefreshTokenRepository interface {
	RefreshTokenSaver
	RefreshTokenFinder
}

type RefreshTokenSaver interface {
	Save(ctx context.Context, token *RefreshToken) error
}

type RefreshTokenFinder interface {
	FindByString(ctx context.Context, str string) (*RefreshToken, error)
}
