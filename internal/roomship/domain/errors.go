package domain

import (
	myErrors "gochat/internal/shared/errors"
)

var (
	ErrInvalidMaxMemberCount      = myErrors.NewBusiness("invalid max member count")
	ErrHandleNotPendingRequest    = myErrors.NewBusiness("handle not pending request")
	ErrContentTooLong             = myErrors.NewBusiness("content too long")
	ErrMemberRequestAlreadyExists = myErrors.NewBusiness("member request already exists")
	ErrIsAlreadyMember            = myErrors.NewBusiness("user is already a member of the room")
	ErrInvalidPassword            = myErrors.NewBusiness("invalid password")
	ErrNotAdmin                   = myErrors.NewBusiness("user is not an admin")
	ErrOwnerCantLeave             = myErrors.NewBusiness("owner can't leave room")
)
