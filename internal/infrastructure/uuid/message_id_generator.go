package uuid

import (
	"gochat/internal/shared/kernel"

	"github.com/google/uuid"
)

var _ kernel.MessageIDGenerator = (*MessageIDGenerator)(nil)

type MessageIDGenerator struct{}

func NewMessageIDGenerator() *MessageIDGenerator {
	return &MessageIDGenerator{}
}

func (m *MessageIDGenerator) Generate() kernel.MessageID {
	id := kernel.MessageID(uuid.Must(uuid.NewV7()).String())
	return id
}
