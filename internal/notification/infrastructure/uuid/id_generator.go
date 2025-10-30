package uuid

import (
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"

	"github.com/google/uuid"
)

var _ application.NoticeIDGenerator = (*NoticeIDGenerator)(nil)

type NoticeIDGenerator struct{}

func NewNoticeIDGenerator() *NoticeIDGenerator {
	return &NoticeIDGenerator{}
}

func (NoticeIDGenerator) Generate() domain.NoticeID {
	return domain.NoticeID(uuid.Must(uuid.NewV7()).String())
}
