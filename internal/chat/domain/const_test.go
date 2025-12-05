package domain_test

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
)

const (
	fixedUserID       kernel.UserID       = "user-123"
	fixedFriendID     kernel.UserID       = "friend-123"
	fixedRoomID       kernel.RoomID       = "room-123"
	fixedFriendshipID domain.FriendshipID = "friendship-123"
	fixedRoomshipID   domain.RoomshipID   = "roomship-123"
)
