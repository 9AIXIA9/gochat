package domain

import "errors"

var (
	ErrEmptyContent = errors.New("message content cannot be empty")
)
