package application_test

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const (
	fixedUserID     kernel.UserID     = "user-123"
	fixedFriendID   kernel.UserID     = "friend-123"
	fixedRoomID     kernel.RoomID     = "room-123"
	fixedMessageID  kernel.MessageID  = "message-123"
	fixedRoomshipID domain.RoomshipID = "roomship-123"
	fixedEventID    event.ID          = "event-123"
	fixedContent    string            = "Hello, how are you?"
)
