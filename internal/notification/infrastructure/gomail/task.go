package gomail

import (
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
)

type task struct {
	mail      *domain.Mail
	email     kernel.Email
	onSuccess func() error
}
