package uuid

import (
	"github.com/google/uuid"
	"gochat/internal/shared/kernel"
)

type IDGenerator struct{}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

func (IDGenerator) Generate() kernel.ID {
	return kernel.ID(uuid.Must(uuid.NewV7()).String())
}
