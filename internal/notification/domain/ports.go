//go:generate mockgen -source=ports.go -destination=./mocks/mock_ports.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type SystemMessageNotifier interface {
	Notify(message *SystemMessage) error
}

type WelcomeEmailNotifier interface {
	NotifyWelcomeEmail(ctx context.Context, email kernel.Email, number kernel.UserNumber) error
}
