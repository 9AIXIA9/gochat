package uuid

import (
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/kernel"

	"github.com/google/uuid"
)

var _ domain.OperationIDGenerator = (*OperationIDGenerator)(nil)

type OperationIDGenerator struct{}

func NewOperationIDGenerator() *OperationIDGenerator {
	return &OperationIDGenerator{}
}

func (OperationIDGenerator) Generate() kernel.OperationID {
	return kernel.OperationID(uuid.Must(uuid.NewV7()).String())
}
