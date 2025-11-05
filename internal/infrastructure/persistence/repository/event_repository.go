package repository

import (
	"context"
	"errors"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/infrastructure/persistence/converter"
	"gochat/internal/infrastructure/persistence/model"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"

	"gorm.io/gorm"
)

var _ event.UnpublishedSaver = (*EventRepository)(nil)
var _ event.UnpublishedLister = (*EventRepository)(nil)
var _ event.PublishedMarker = (*EventRepository)(nil)
var _ event.DeadLetterSaver = (*EventRepository)(nil)

type EventRepository struct {
	db                 *gorm.DB
	inner              *gormutils.Repository[model.Event, event.StandardEvent]
	modelConverter     gormutils.GenericModelConverter[*model.Event, *event.StandardEvent]
	interfaceConverter *converter.EventInterfaceConverter
}

func NewEventRepository(db *gorm.DB, modelConverter gormutils.GenericModelConverter[*model.Event, *event.StandardEvent]) *EventRepository {
	return &EventRepository{
		db:                 db,
		inner:              gormutils.NewRepository(db, modelConverter),
		modelConverter:     modelConverter,
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

func (repo *EventRepository) MarkPublished(ctx context.Context, ID event.ID) error {

	if err := repo.inner.Delete(ctx, "id = ?", ID); !errors.Is(err, myErrors.ErrNotFound) {
		return err
	}
	return nil
}

func (repo *EventRepository) SaveDeadLetter(ctx context.Context, event event.Event, reason error) error {
	if event == nil {
		return nil
	}
	DeadLetterModel := model.NewDeadLetter(repo.modelConverter.ToModel(repo.interfaceConverter.ToStandard(event)), reason)
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(DeadLetterModel).Error)
}
