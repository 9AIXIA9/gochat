package application

import (
	"context"
	"gochat/internal/shared/kernel"
)

type UserExister interface {
	ExistsByID(ctx context.Context, userID kernel.UserID) (bool, error)
}
