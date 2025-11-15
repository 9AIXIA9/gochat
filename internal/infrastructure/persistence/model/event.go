package model

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

type Event struct {
	ID          event.ID    `gorm:"primaryKey;type:char(36)"`
	AggregateID kernel.ID   `gorm:"type:char(36);not null;index"`
	Topic       event.Topic `gorm:"type:varchar(100);not null;index"`
	Payload     []byte

	CreatedAt time.Time
}

func (e *Event) TableName() string {
	return "gochat.unpublished_events"
}
