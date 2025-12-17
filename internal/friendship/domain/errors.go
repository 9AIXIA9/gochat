package domain

import (
	myErrors "gochat/internal/shared/errors"
)

var (
	ErrAddYourselfAsFriend         = myErrors.NewBusiness("add yourself as a friend")
	ErrAlreadyBeenFriends          = myErrors.NewBusiness("already been friends")
	ErrFriendRequestExists         = myErrors.NewBusiness("friend request already exists")
	ErrFriendRequestNotForUser     = myErrors.NewBusiness("friend request not for user")
	ErrFriendRequestHasBeenHandled = myErrors.NewBusiness("friend request has been handled")
	ErrFriendRequestContentTooLong = myErrors.NewBusiness("friend request content too long")
	ErrFriendRequestNotAgreed      = myErrors.NewBusiness("friend request is not agreed")
)
