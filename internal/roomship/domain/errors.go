package domain

import "errors"

var (
	ErrInvalidMaxMemberCount   = errors.New("invalid max member count")
	ErrHandleNotPendingRequest = errors.New("handle not pending request")
)
