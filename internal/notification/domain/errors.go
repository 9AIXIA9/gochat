package domain

import (
	"errors"
)

var (
	ErrNotificationNoRecipient = errors.New("notification has no recipient")
)
