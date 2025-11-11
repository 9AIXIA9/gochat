package application

import (
	"context"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
)

type UserCreatedEmailNotifier interface {
	AddUserCreatedEmail(ctx context.Context, email kernel.Email, number domain.UserNumber) error
}

type MessageNotifier interface {
	Notify(recipient kernel.UserID, message *domain.Message) error
}
