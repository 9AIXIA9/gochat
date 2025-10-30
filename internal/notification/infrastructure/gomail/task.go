package gomail

import (
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
)

type task struct {
	notice    *domain.Notice
	email     kernel.Email
	onSuccess func() error
}
