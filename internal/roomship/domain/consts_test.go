package domain_test

import (
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const (
	maxContentLength = 100

	timeTolerance                                   = 150 * time.Millisecond
	fixedUserID            kernel.UserID            = "user-0001"
	fixedRoomID            kernel.RoomID            = "room-0001"
	fixedOperationID       kernel.OperationID       = "operation-0001"
	fixedOperatorID        kernel.UserID            = "operator-1000"
	fixedEventID           event.ID                 = "event-123"
	fixedContent                                    = "hello"
	fixedRoomNumber        domain.RoomNumber        = "123456789"
	fixedPassword          domain.Password          = "12345"
	fixedPasswordEncrypted domain.PasswordEncrypted = "encrypted-12345"
	fixedRoomshipID        domain.RoomshipID        = "roomship-0001"
)
