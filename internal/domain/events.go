package domain

import "time"

// BaseEvent 基础事件结构
type BaseEvent struct {
	id         EventID
	occurredAt time.Time
}

func (b BaseEvent) EventID() EventID {
	return b.id
}

func (b BaseEvent) OccurredOn() time.Time {
	return b.occurredAt
}
