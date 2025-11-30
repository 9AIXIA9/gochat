package domain

import "errors"

var (
	ErrAddYourselfAsFriend     = errors.New("add yourself as a friend")
	ErrAlreadyBeenFriends      = errors.New("already been friends")
	ErrFriendRequestExists     = errors.New("friend request already exists")
	ErrFriendRequestNotForUser = errors.New("friend request not for user")
)
