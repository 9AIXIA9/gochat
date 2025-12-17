package domain

import (
	myErrors "gochat/internal/shared/errors"
)

var (
	ErrEmptyMessageContent = myErrors.NewBusiness("message content cannot be empty")
	ErrNotFriends          = myErrors.NewBusiness("users are not friends")
	ErrNotMember           = myErrors.NewBusiness("user is not a member of the room")
)
