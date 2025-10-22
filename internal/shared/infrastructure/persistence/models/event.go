package models

import (
	kernelEvent "gochat/internal/shared/kernel/event"
	"time"
)

type Event struct {
	ID         kernelEvent.ID `gorm:"primaryKey;type:char(36)"`
	Topic      kernelEvent.Topic
	OccurredAt time.Time
	Payload    []byte
}

func (e *Event) TableName() string {
	return "gochat.unpublished_events"
}
