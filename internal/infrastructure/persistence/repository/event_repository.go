package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/infrastructure/persistence/model"
	"gochat/internal/shared/event"
)

var _ event.Repository = (*EventRepository)(nil)

type EventRepository struct {
	unitOfWork *gormutils.UnitOfWork
}

func NewEventRepository(unitOfWork *gormutils.UnitOfWork) *EventRepository {
	return &EventRepository{
		unitOfWork: unitOfWork,
	}
}

func (repo *EventRepository) Save(ctx context.Context, event event.Event) error {
	ev := &model.Event{
		ID:          event.ID(),
		AggregateID: event.AggregateID(),
		Topic:       event.Topic(),
		Published:   false,
		Payload:     event.Payload(),
		CreatedAt:   event.OccurredAt(),
	}
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Create(ev).Error)
}

func (repo *EventRepository) Saves(ctx context.Context, events []event.Event) error {
	if len(events) == 0 {
		return nil
	}
	models := make([]*model.Event, 0, len(events))
	for _, e := range events {
		models = append(models, &model.Event{
			ID:          e.ID(),
			AggregateID: e.AggregateID(),
			Topic:       e.Topic(),
			Published:   false,
			Payload:     e.Payload(),
			CreatedAt:   e.OccurredAt(),
		})
	}
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Create(&models).Error)
}

func (repo *EventRepository) UnpublishedList(ctx context.Context) ([]event.Event, error) {
	var models []*model.Event

	if err := repo.unitOfWork.DB(ctx).WithContext(ctx).Where("published = false").Find(&models).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	events := make([]event.Event, 0, len(models))
	for _, m := range models {
		events = append(events, event.NewStandardEvent(m.ID, m.AggregateID, m.CreatedAt, m.Topic, m.Payload))
	}
	return events, nil
}

func (repo *EventRepository) MarkPublished(ctx context.Context, ID event.ID) error {
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Model(&model.Event{}).Where("id = ?", ID).Update("published", true).Error)
}

func (repo *EventRepository) SaveDeadLetter(ctx context.Context, event event.Event, reason error) error {
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Create(&model.DeadLetter{
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
