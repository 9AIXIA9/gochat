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
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(repo.toModel(notification)).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		if err := repo.eventRepo.CreateUnpublishedEvents(gormutils.SetTransaction(ctx, tx), notification.GetEvents()); err != nil {
			return gormutils.TranslateError(err)
		}

		return nil
	})
}

func (repo *NotificationRepository) UpdateStateByID(ctx context.Context, id kernel.MessageID, state domain.NotificationState) error {
	if err := repo.db.WithContext(ctx).Model(&model.Notification{}).
		Where("id = ?", id).
		Update("state", state).Error; err != nil {
		return gormutils.TranslateError(err)
	}
	return nil
}

func (repo *NotificationRepository) toModel(notification *domain.Notification) *model.Notification {
	return &model.Notification{
		ID:          notification.ID(),
		RecipientID: notification.RecipientID(),
		State:       notification.State(),
		RawPayload:  notification.RawPayload(),
	}
}
