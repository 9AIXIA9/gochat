package application

import (
	"context"
	"gochat/internal/shared/kernel"
)

type WelcomeEmailNotifier interface {
	NotifyWelcomeEmail(ctx context.Context, email kernel.Email, number kernel.UserNumber) error
}
