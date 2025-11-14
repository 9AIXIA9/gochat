package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/infrastructure/persistence/model"
	"gochat/internal/shared/event"

	"gorm.io/gorm"
)

var _ event.Repository = (*EventRepository)(nil)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (repo *EventRepository) Saves(ctx context.Context, events []event.Event) error {
	models := make([]*model.Event, 0, len(events))
	for _, e := range events {
		models = append(models, &model.Event{
			ID:          e.ID(),
			AggregateID: e.AggregateID(),
			Topic:       e.Topic(),
			Payload:     e.Payload(),
			CreatedAt:   e.OccurredAt(),
		})
	}
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(&models).Error)
}

func (repo *EventRepository) UnpublishedList(ctx context.Context) ([]event.Event, error) {
	var models []*model.Event
	if err := repo.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	events := make([]event.Event, 0, len(models))
	for _, m := range models {
		events = append(events, event.NewStandardEvent(m.ID, m.AggregateID, m.CreatedAt, m.Topic, m.Payload))
	}
	return events, nil
}

func (repo *EventRepository) MarkPublished(ctx context.Context, ID event.ID) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Delete(&model.Event{}, "id = ?", ID).Error)
}

func (repo *EventRepository) SaveDeadLetter(ctx context.Context, event event.Event, reason error) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(&model.DeadLetter{
		Event: &model.Event{
			ID:          event.ID(),
			AggregateID: event.AggregateID(),
			Topic:       event.Topic(),
			Payload:     event.Payload(),
			CreatedAt:   event.OccurredAt(),
		},
		Reason: reason.Error(),
	}).Error)
}
