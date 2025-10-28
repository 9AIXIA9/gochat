package model

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	ID          event.ID `gorm:"primaryKey;type:char(36)"`
	AggregateID kernel.ID
	Topic       event.Topic
	OccurredAt  time.Time
	Payload     []byte
}

func (e *Event) TableName() string {
	return "gochat.unpublished_events"
}
