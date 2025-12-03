package domain

import "errors"

var (
	ErrNotBelongToRoom = errors.New("user does not belong to the room")
)
