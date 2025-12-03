package domain

import "errors"

var (
	ErrEmptyContent   = errors.New("message content cannot be empty")
	ErrNotUndelivered = errors.New("message is not undelivered")
)
