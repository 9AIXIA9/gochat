package repository

import (
	"context"
	"encoding/json"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/infrastructure/persistence/model"
	"gochat/internal/shared/event"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	for _, ev := range evs {
		setEventTrace(ctx, ev)
	}

	tx := gormutils.GetTransaction(ctx)
	if tx == nil {
		tx = repo.db
	}

	evModel, err := repo.toModels(evs)
	if err != nil {
		return err
	}

	if err := gormutils.TranslateError(tx.WithContext(ctx).Create(evModel).Error); err != nil {
		return err
	}
	return nil
}

func (repo *EventRepository) ListUnpublishedEvents(ctx context.Context, lease time.Duration) ([]event.Event, error) {
	var models []*model.Event
	now := time.Now().UTC()
	leaseUntil := now.Add(lease)

	if err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("processing_until <= ?", now).
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Limit(1000).
			Find(&models).
			Error; err != nil {
			return gormutils.TranslateError(err)
		}

		if len(models) == 0 {
			return nil
		}

		ids := make([]event.ID, 0, len(models))
		for _, m := range models {
			ids = append(ids, m.ID)
		}

		if err := tx.Model(&model.Event{}).
			Where("id IN ?", ids).
			Update("processing_until", leaseUntil).
			Error; err != nil {
			return gormutils.TranslateError(err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return repo.toEvents(models)
}

func (repo *EventRepository) MarkAsPublished(ctx context.Context, ID event.ID) error {
	return gormutils.TranslateError(
		repo.db.WithContext(ctx).
			Where("id = ?", ID).
			Delete(&model.Event{}).
			Error,
	)
}

func (repo *EventRepository) CreateDeadLetter(ctx context.Context, e event.Event, reason error) error {
	if e == nil {
		return nil
	}
	deadLetter, err := repo.toDeadLetter(e, reason)
	if err != nil {
		return err
	}

	return gormutils.TranslateError(repo.db.
		WithContext(ctx).
		Create(deadLetter).
		Error)
}

func (repo *EventRepository) toModel(e event.Event) (*model.Event, error) {
	headers := e.Headers()
	data, err := json.Marshal(&headers)
	if err != nil {
		return nil, err
	}

	return &model.Event{
		ID:              e.ID(),
		AggregateID:     e.AggregateID(),
		Topic:           e.Topic(),
		ProcessingUntil: time.Now().UTC(),
		Payload:         e.Payload(),
		Headers:         data,
		CreatedAt:       e.OccurredAt(),
	}, nil
}

func (repo *EventRepository) toModels(evs []event.Event) ([]*model.Event, error) {
	if len(evs) == 0 {
		return nil, nil
	}
	models := make([]*model.Event, 0, len(evs))
	for _, e := range evs {
		modelEv, err := repo.toModel(e)
		if err != nil {
			return nil, err
		}
		models = append(models, modelEv)
	}
	return models, nil
}

func (repo *EventRepository) toEvent(model *model.Event) (event.Event, error) {
	if model == nil {
		return nil, nil
	}

	if model.Headers == nil {
		return event.LoadStandardEvent(
			model.ID,
			model.AggregateID,
			model.CreatedAt,
			model.Topic,
			model.Payload,
			nil,
		), nil
	}

	var headers map[string]string
	if err := json.Unmarshal(model.Headers, &headers); err != nil {
		return nil, err
	}

	return event.LoadStandardEvent(
		model.ID,
		model.AggregateID,
		model.CreatedAt,
		model.Topic,
		model.Payload,
		headers,
	), nil
}

func (repo *EventRepository) toEvents(models []*model.Event) ([]event.Event, error) {
	if len(models) == 0 {
		return nil, nil
	}
	evs := make([]event.Event, 0, len(models))
	for _, m := range models {
		ev, err := repo.toEvent(m)
		if err != nil {
			return nil, err
		}
		evs = append(evs, ev)
	}
	return evs, nil
}

func (repo *EventRepository) toDeadLetter(e event.Event, reason error) (*model.DeadLetter, error) {
	evModel, err := repo.toModel(e)
	if err != nil {
		return nil, err
	}
	return &model.DeadLetter{
		Event:  evModel,
		Reason: reason.Error(),
	}, nil
}

func setEventTrace(ctx context.Context, ev event.Event) {
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	ev.AddHeaders(carrier)
}
