package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.NotificationRepository = (*NotificationRepository)(nil)

type NotificationRepository struct {
	db        *gorm.DB
	eventRepo event.Repository
}

func NewNotificationRepository(db *gorm.DB, eventRepo event.Repository) *NotificationRepository {
	return &NotificationRepository{
		db:        db,
		eventRepo: eventRepo,
	}
}

func (repo *NotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	return repo.db.WithContext(ctx).Create(repo.toModel(notification)).Error
}

func (repo *NotificationRepository) Delete(ctx context.Context, id kernel.MessageID) error {
	if err := repo.db.WithContext(ctx).Model(&model.Notification{}).
		Where("id = ?", id).Error; err != nil {
		return gormutils.TranslateError(err)
	}
	return nil
}

func (repo *NotificationRepository) toModel(notification *domain.Notification) *model.Notification {
	return &model.Notification{
		ID:          notification.ID(),
		RecipientID: notification.RecipientID(),
		RawPayload:  notification.RawPayload(),
	}
}
