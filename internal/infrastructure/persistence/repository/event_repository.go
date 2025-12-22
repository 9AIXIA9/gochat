package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/infrastructure/persistence/model"
	"gochat/internal/shared/event"
	"time"

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

func (repo *EventRepository) CreateUnpublishedEvents(ctx context.Context, evs []event.Event) error {
	if len(evs) == 0 {
		return nil
	}

	tx := gormutils.GetTransaction(ctx)
	if tx == nil {
		tx = repo.db
	}

	if err := gormutils.TranslateError(tx.WithContext(ctx).Create(repo.toModels(evs)).Error); err != nil {
		return err
	}
	return nil
}

func (repo *EventRepository) ListUnpublishedEvents(ctx context.Context) ([]event.Event, error) {
	var models []*model.Event

	if err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("published = ? AND processing = ?", false, false).Find(&models).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		if len(models) == 0 {
			return nil
		}

		ids := make([]event.ID, 0, len(models))
		for _, m := range models {
			ids = append(ids, m.ID)
		}

		if err := tx.Model(&model.Event{}).Where("id IN ?", ids).Updates(map[string]interface{}{
			"processing": true,
		}).Error; err != nil {
			return gormutils.TranslateError(err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return repo.toEvents(models), nil
}

func (repo *EventRepository) MarkAsPublishedAndNotProcessing(ctx context.Context, ID event.ID) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Model(&model.Event{}).Where("id = ?", ID).Updates(map[string]interface{}{
		"published":    true,
		"processing":   false,
		"published_at": time.Now().UTC(),
	}).Error)
}

func (repo *EventRepository) MarkAsNotProcessing(ctx context.Context, ID event.ID) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Model(&model.Event{}).Where("id = ?", ID).Update("processing", false).Error)
}

func (repo *EventRepository) CreateDeadLetter(ctx context.Context, e event.Event, reason error) error {
	if e == nil {
		return nil
	}
	return gormutils.TranslateError(repo.db.
		WithContext(ctx).
		Create(repo.toDeadLetter(e, reason)).
		Error)
}

func (repo *EventRepository) toModel(e event.Event) *model.Event {
	return &model.Event{
		ID:          e.ID(),
		AggregateID: e.AggregateID(),
		Topic:       e.Topic(),
		Published:   false,
		Processing:  false,
		Payload:     e.Payload(),
		CreatedAt:   e.OccurredAt(),
	}
}

func (repo *EventRepository) toModels(evs []event.Event) []*model.Event {
	if len(evs) == 0 {
		return nil
	}
	models := make([]*model.Event, 0, len(evs))
	for _, e := range evs {
		models = append(models, repo.toModel(e))
	}
	return models
}

func (repo *EventRepository) toEvent(model *model.Event) event.Event {
	if model == nil {
		return nil
	}
	return event.LoadStandardEvent(
		model.ID,
		model.AggregateID,
		model.CreatedAt,
		model.Topic,
		model.Payload,
	)
}

func (repo *EventRepository) toEvents(models []*model.Event) []event.Event {
	if len(models) == 0 {
		return nil
	}
	evs := make([]event.Event, 0, len(models))
	for _, m := range models {
		evs = append(evs, repo.toEvent(m))
	}
	return evs
}

func (repo *EventRepository) toDeadLetter(e event.Event, reason error) *model.DeadLetter {
	return &model.DeadLetter{
		Event:  repo.toModel(e),
		Reason: reason.Error(),
	}
}
