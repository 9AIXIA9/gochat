package application

import (
	"context"
	"gochat/internal/shared/kernel"
)

type UserRepository interface {
	UserIDSaver
}

type UserIDSaver interface {
	SaveID(ctx context.Context, userID kernel.UserID) error
}
