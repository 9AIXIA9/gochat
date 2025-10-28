package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/infrastructure/persistence/converter"
	"gochat/internal/infrastructure/persistence/model"
	"gochat/internal/shared/event"

	"gorm.io/gorm"
)

var _ event.Repository = (*EventRepository)(nil)

type EventRepository struct {
	db                 *gorm.DB
	inner              *gormutils.Repository[model.Event, event.StandardEvent]
	interfaceConverter *converter.EventInterfaceConverter
}

func NewEventRepository(db *gorm.DB, modelConverter gormutils.GenericModelConverter[*model.Event, *event.StandardEvent]) *EventRepository {
	return &EventRepository{
		db:                 db,
		inner:              gormutils.NewRepository(db, modelConverter),
		interfaceConverter: &converter.EventInterfaceConverter{},
	}
}

func (repo *EventRepository) Saves(ctx context.Context, events []event.Event) error {
	return repo.inner.Saves(ctx, repo.interfaceConverter.ToStandards(events))
}

func (repo *EventRepository) UnpublishedList(ctx context.Context) ([]event.Event, error) {
	events, err := repo.inner.Finds(ctx)
	if err != nil {
		return nil, err
	}
	return repo.interfaceConverter.ToInterfaces(events), nil
}

func (repo *EventRepository) MarkPublished(ctx context.Context, IDs []event.ID) error {
	return repo.inner.Delete(ctx, "id IN ?", IDs)
}
