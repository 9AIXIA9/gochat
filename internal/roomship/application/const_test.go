package application_test

import (
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const (
	defaultMemberCount = 20

	fixedRoomshipID        domain.RoomshipID        = "roomship-123"
	fixedUserID            kernel.UserID            = "user-123"
	fixedOwnerID           kernel.UserID            = "owner-123"
	fixedApplicantID       kernel.UserID            = "applicant-123"
	fixedRoomID            kernel.RoomID            = "room-123"
	fixedOperationID       kernel.OperationID       = "operation-123"
	fixedEventID           event.ID                 = "event-123"
	fixedContent                                    = "hello"
	fixedOperatorID        kernel.UserID            = "operator-123"
	fixedMaxMemberCount                             = 50
	fixedPassword          domain.Password          = "secure-password-123"
	fixedPasswordEncrypted domain.PasswordEncrypted = "encrypted-password-123"
	fixedRoomNumber        domain.RoomNumber        = "1001"
)
