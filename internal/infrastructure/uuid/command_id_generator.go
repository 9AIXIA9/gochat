package uuid

import (
	"gochat/internal/shared/command"

	"github.com/google/uuid"
)

var _ command.IDGenerator = (*CommandIDGenerator)(nil)

type CommandIDGenerator struct{}

func NewCommandIDGenerator() *CommandIDGenerator {
	return &CommandIDGenerator{}
}

func (CommandIDGenerator) Generate() command.ID {
	return command.ID(uuid.Must(uuid.NewV7()).String())
}
