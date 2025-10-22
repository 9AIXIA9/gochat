package repository

import (
	"context"

	gormutils "gochat/internal/shared/infrastructure/gorm"
	"gochat/internal/shared/infrastructure/persistence/converters"
	"gochat/internal/shared/infrastructure/persistence/models"
	kernelEvent "gochat/internal/shared/kernel/event"
	"gorm.io/gorm"
)

type PublisherRepository[SpecificEvent kernelEvent.Event] struct {
	db        *gorm.DB
	converter *converters.ModelToSpecificEventConverter[SpecificEvent]
}

func NewPublisherRepository[SpecificEvent kernelEvent.Event](db *gorm.DB, converter *converters.ModelToSpecificEventConverter[SpecificEvent]) *PublisherRepository[SpecificEvent] {
	return &PublisherRepository[SpecificEvent]{
		db:        db,
		converter: converter,
	}
}

func (r *PublisherRepository[SpecificEvent]) ListUnpublishedEvent(ctx context.Context) ([]SpecificEvent, error) {
	var gormEvents []*models.Event
	if err := r.db.WithContext(ctx).
		Order("occurred_at ASC").
		Find(&gormEvents).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return r.converter.ToModels(gormEvents)
}

func (r *PublisherRepository[SpecificEvent]) MarkAsPublished(ctx context.Context, IDs []kernelEvent.ID) error {
	if len(IDs) == 0 {
		return nil
	}
	tx := r.db.WithContext(ctx).
		Where("id IN ?", IDs).
		Delete(&models.Event{})
	return gormutils.TranslateError(tx.Error)
}
