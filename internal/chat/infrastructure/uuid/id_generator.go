package uuid

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"

	"github.com/google/uuid"
)

var _ application.MessageIDGenerator = (*MessageIDGenerator)(nil)

type MessageIDGenerator struct{}

func NewMessageIDGenerator() *MessageIDGenerator {
	return &MessageIDGenerator{}
}

func (MessageIDGenerator) Generate() domain.MessageID {
	return domain.MessageID(uuid.Must(uuid.NewV7()).String())
}
