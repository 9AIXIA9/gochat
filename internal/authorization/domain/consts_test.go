package domain_test

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const (
	fixedUserID            kernel.UserID            = "user-123"
	fixedUserNumber        kernel.UserNumber        = "10001"
	fixedPassword          domain.Password          = "plain-pass"
	fixedPasswordEncrypted domain.PasswordEncrypted = "encrypted-pass"
	fixedEventID           event.ID                 = "event-999"
	fixedRefreshToken      domain.RefreshToken      = "refresh-token-abc"
	fixedAccessToken       domain.AccessToken       = "access-token-xyz"
	fixedEmail             kernel.Email             = "test@example.com"
	fixedRefreshCount                               = 3

	timeTolerance    = 150 * time.Millisecond
	validityDuration = 7 * 24 * time.Hour // mirrored business rule
	extendDuration   = 24 * time.Hour     // mirrored business rule
	maxRefreshCount  = 7 * 24             // mirrored business rule
)
