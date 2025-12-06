package application_test

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const (
	fixedUserID            kernel.UserID            = "user-123"
	fixedUserNumber        kernel.UserNumber        = "12345"
	fixedPassword          domain.Password          = "password"
	fixedEncryptedPassword domain.PasswordEncrypted = "encrypted-password"
	fixedEmail             kernel.Email             = "email@demo.com"
	fixedEventID           event.ID                 = "event-12345"
	fixedRefreshToken      domain.RefreshToken      = "refresh"
	fixedAccessToken       domain.AccessToken       = "access"
	fixedRefreshCount                               = 3
)
