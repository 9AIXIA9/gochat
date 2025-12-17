package domain

import (
	myErrors "gochat/internal/shared/errors"
)

var (
	ErrEmptyContent = myErrors.NewBusiness("message content cannot be empty")
)
