package domain

import "time"

// Event - 领域事件接口
type Event interface {
	EventType() string
	OccurredOn() time.Time
}

type RoomCreatedEvent struct {
}

type RoomJoinedEvent struct {
}

type RoomLeftEvent struct {
}
