package domain_test

import (
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const (
	fixedUserID       kernel.UserID       = "user-123"
	fixedToUserID     kernel.UserID       = "user-9999"
	fixedRequestID    kernel.OperationID  = "friend-request-456"
	fixedEventID      event.ID            = "ev-6789"
	fixedContent      string              = "Hello, let's be friends!"
	fixedFriendshipID domain.FriendshipID = "friendship-321"
	timeTolerance                         = 150 * time.Millisecond
	maxContentLength                      = 100
)
