package domain_test

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const (
	fixedUserID         kernel.UserID    = "user-123"
	fixedMemberID1      kernel.UserID    = "member-1"
	fixedMemberID2      kernel.UserID    = "member-2"
	fixedRoomID         kernel.RoomID    = "room-abc"
	fixedMessageID      kernel.MessageID = "msg-123"
	fixedEventID        event.ID         = "event-1234"
	fixedMessageContent                  = "hello"
	timeTolerance                        = 500 * time.Millisecond
)
