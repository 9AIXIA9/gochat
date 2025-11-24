package domain

import "gochat/internal/shared/kernel"

type MessageIDGenerator interface {
	Generate() kernel.MessageID
}
