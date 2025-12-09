package application_test

import (
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const (
	fixedUserID       kernel.UserID       = "mock-user-123"
	fixedToID         kernel.UserID       = "to-user-456"
	fixedOperationID  kernel.OperationID  = "operation-789"
	fixedContent                          = "Let's be friends!"
	fixedEventID      event.ID            = "event-0001"
	fixedFriendshipID domain.FriendshipID = "friend-23333"
)
