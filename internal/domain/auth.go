package domain

import (
	"context"
)

type AuthUsecase interface {
	ParseToken(ctx context.Context, tokenStr string) bool
}
