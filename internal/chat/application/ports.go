package application

import "gochat/internal/shared/kernel"

type MessageIDGenerator interface {
	Generate() kernel.MessageID
}
