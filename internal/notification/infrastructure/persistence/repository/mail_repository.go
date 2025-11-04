package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"

	"gorm.io/gorm"
)

var _ application.MailSaver = (*MailRepository)(nil)

type MailRepository struct {
	innerRepository *gormutils.Repository[model.Mail, domain.Mail]
	db              *gorm.DB
}

func NewMailRepository(db *gorm.DB, converter gormutils.GenericModelConverter[*model.Mail, *domain.Mail]) *MailRepository {
	return &MailRepository{innerRepository: gormutils.NewRepository(db, converter), db: db}
}

func (repo *MailRepository) Save(ctx context.Context, mail *domain.Mail) error {
	return repo.innerRepository.Save(ctx, mail)
}
