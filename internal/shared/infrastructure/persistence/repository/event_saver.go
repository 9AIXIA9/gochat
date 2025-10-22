package repository

import (
	"context"
	gormutils "gochat/internal/shared/infrastructure/gorm"
	"gochat/internal/shared/infrastructure/persistence/converters"
	kernelEvent "gochat/internal/shared/kernel/event"
	"gorm.io/gorm"
)

var _ kernelEvent.Saver = (*EventSaver)(nil)

type EventSaver struct {
	db        *gorm.DB
	converter *converters.EventToModelConverter
}

func NewEventSaver(db *gorm.DB, converter *converters.EventToModelConverter) *EventSaver {
	return &EventSaver{
		db:        db,
		converter: converter,
	}
}

func (s *EventSaver) Save(ctx context.Context, events []kernelEvent.Event) error {
	if len(events) == 0 {
		return nil
	}
	gormEvents, err := s.converter.ToModels(events)
	if err != nil {
		return err
	}
	return gormutils.TranslateError(s.db.WithContext(ctx).Create(gormEvents).Error)
}
