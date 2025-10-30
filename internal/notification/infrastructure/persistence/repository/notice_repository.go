package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"

	"gorm.io/gorm"
)

var _ application.NoticeRepository = (*NoticeRepository)(nil)

type NoticeRepository struct {
	innerRepository *gormutils.Repository[model.Notice, domain.Notice]
	db              *gorm.DB
}

func NewNoticeRepository(db *gorm.DB, converter gormutils.GenericModelConverter[*model.Notice, *domain.Notice]) *NoticeRepository {
	return &NoticeRepository{innerRepository: gormutils.NewRepository(db, converter), db: db}
}

func (repo *NoticeRepository) Save(ctx context.Context, notice *domain.Notice) error {
	return repo.innerRepository.Save(ctx, notice)
}
