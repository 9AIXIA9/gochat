package application

import (
	"gochat/internal/chat/domain"
)

type MessageIDGenerator interface {
	Generate() domain.MessageID
}
