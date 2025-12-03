package domain_test

import (
	"gochat/internal/shared/kernel"
	"time"
)

const (
	timeTolerance                     = 500 * time.Millisecond
	fixedMessageID   kernel.MessageID = "fixed-message-id"
	fixedRecipientID kernel.UserID    = "fixed-recipient-id"
	fixedContent     string           = "This is a system message."
)
