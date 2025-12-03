package repository

import (
	"context"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.SystemMessageRepository = (*SystemMessageRepository)(nil)

type SystemMessageRepository struct {
	db *gorm.DB
}

func NewSystemMessageRepository(db *gorm.DB) *SystemMessageRepository {
	return &SystemMessageRepository{db: db}
}

func (repo *SystemMessageRepository) Create(ctx context.Context, message *domain.SystemMessage) error {
	//TODO implement me
	panic("implement me")
}

func (repo *SystemMessageRepository) Update(ctx context.Context, message *domain.SystemMessage) error {
	//TODO implement me
	panic("implement me")
}

func (repo *SystemMessageRepository) Updates(ctx context.Context, messages []*domain.SystemMessage) error {
	//TODO implement me
	panic("implement me")
}

func (repo *SystemMessageRepository) FindByID(ctx context.Context, messageID kernel.MessageID) (*domain.SystemMessage, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *SystemMessageRepository) FindsByState(ctx context.Context, userID kernel.UserID, state domain.MessageState) ([]*domain.SystemMessage, error) {
	//TODO implement me
	panic("implement me")
}
