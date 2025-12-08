package domain

import "errors"

var (
	ErrInvalidMaxMemberCount      = errors.New("invalid max member count")
	ErrHandleNotPendingRequest    = errors.New("handle not pending request")
	ErrContentTooLong             = errors.New("content too long")
	ErrMemberRequestAlreadyExists = errors.New("member request already exists")
	ErrIsAlreadyMember            = errors.New("user is already a member of the room")
	ErrInvalidPassword            = errors.New("invalid password")
	ErrNotAdmin                   = errors.New("user is not an admin")
)
