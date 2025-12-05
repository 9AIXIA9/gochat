package domain

import "errors"

var (
	ErrNotFriends = errors.New("users are not friends")
	ErrNotMember  = errors.New("user is not a member of the room")
)
