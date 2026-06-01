package uuid

import (
	"gochat/internal/shared/command"

	"github.com/google/uuid"
)

var _ command.ReceiptIDGenerator = (*CommandReceiptIDGenerator)(nil)

type CommandReceiptIDGenerator struct{}

func NewCommandReceiptIDGenerator() *CommandReceiptIDGenerator {
	return &CommandReceiptIDGenerator{}
}

func (CommandReceiptIDGenerator) Generate() command.ReceiptID {
	return command.ReceiptID(uuid.Must(uuid.NewV7()).String())
}
