package repository

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"

	"gorm.io/gorm"
)

var _ application.MessagesSaver = (*MessageRepository)(nil)

type MessageRepository struct {
	innerRepository *gormutils.Repository[model.Message, domain.Message]
	db              *gorm.DB
}

func NewMessageRepository(db *gorm.DB, converter gormutils.GenericModelConverter[*model.Message, *domain.Message]) *MessageRepository {
	return &MessageRepository{innerRepository: gormutils.NewRepository(db, converter), db: db}
}

func (repo *MessageRepository) Save(ctx context.Context, messages []*domain.Message) error {
	return repo.innerRepository.Saves(ctx, messages)
}
