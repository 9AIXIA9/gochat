package uuid

import (
	"gochat/internal/gateway/core"

	"github.com/google/uuid"
)

var _ core.SessionIDGenerator = (*SessionIDGenerator)(nil)

type SessionIDGenerator struct{}

func NewSessionIDGenerator() *SessionIDGenerator {
	return &SessionIDGenerator{}
}

func (SessionIDGenerator) Generate() core.SessionID {
	return core.SessionID(uuid.Must(uuid.NewV7()).String())
}
