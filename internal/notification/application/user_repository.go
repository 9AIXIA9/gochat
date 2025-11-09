package application

import (
	"context"
	"gochat/internal/shared/kernel"
)

type UserRepository interface {
	UserEmailSaver
}

type UserEmailSaver interface {
	SaveEmail(ctx context.Context, userID kernel.UserID, email kernel.Email) error
}
