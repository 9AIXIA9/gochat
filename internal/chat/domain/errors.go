package domain

import (
	myErrors "gochat/internal/shared/errors"
)

var (
	ErrNotFriends = myErrors.NewBusiness("users are not friends")
	ErrNotMember  = myErrors.NewBusiness("user is not a member of the room")
)
