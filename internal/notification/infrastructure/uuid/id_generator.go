package uuid

import (
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"

	"github.com/google/uuid"
)

var _ application.MailIDGenerator = (*MailIDGenerator)(nil)

type MailIDGenerator struct{}

func NewMailIDGenerator() *MailIDGenerator {
	return &MailIDGenerator{}
}

func (MailIDGenerator) Generate() domain.MailID {
	return domain.MailID(uuid.Must(uuid.NewV7()).String())
}
