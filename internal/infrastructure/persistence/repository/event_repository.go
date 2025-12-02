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

func (repo *EventRepository) CreateUnpublishedEvent(ctx context.Context, e event.Event) error {
	if e == nil {
		return nil
	}
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(repo.toModel(e)).Error)
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

	if err := repo.db.WithContext(ctx).Where("published = false").Find(&models).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toEvents(models), nil
}

func (repo *EventRepository) MarkAsPublished(ctx context.Context, ID event.ID) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Model(&model.Event{}).Where("id = ?", ID).Update("published", true).Error)
}

func (repo *EventRepository) CreateDeadLetter(ctx context.Context, e event.Event, reason error) error {
	if e == nil {
		return nil
	}
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(repo.toDeadLetter(e, reason)).Error)
}

func (repo *EventRepository) toModel(e event.Event) *model.Event {
	return &model.Event{
		ID:          e.ID(),
		AggregateID: e.AggregateID(),
		Topic:       e.Topic(),
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
		models = append(models, &model.Event{
			ID:          e.ID(),
			AggregateID: e.AggregateID(),
			Topic:       e.Topic(),
			Published:   false,
			Payload:     e.Payload(),
			CreatedAt:   e.OccurredAt(),
		})
	}
	return models
}

func (repo *EventRepository) toDeadLetter(e event.Event, reason error) *model.DeadLetter {
	return &model.DeadLetter{
		Event: &model.Event{
			ID:          e.ID(),
			AggregateID: e.AggregateID(),
			Topic:       e.Topic(),
			Payload:     e.Payload(),
			CreatedAt:   e.OccurredAt(),
		},
		Reason: reason.Error(),
	}
}

func (repo *EventRepository) toEvents(models []*model.Event) []event.Event {
	if len(models) == 0 {
		return nil
	}
	events := make([]event.Event, 0, len(models))
	for _, m := range models {
		events = append(events, event.LoadStandardEvent(m.ID, m.AggregateID, m.CreatedAt, m.Topic, m.Payload))
	}
	return events
}
