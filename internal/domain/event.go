package domain

import (
	"time"
)

type EventTopic string

//todo 完成接口

type Publisher interface {
	Close() error
}

type Subscriber interface {
	Close() error
}

type Event interface {
	EventTopic() EventTopic
	OccurredOn() time.Time
}
type UserCreatedEvent struct {
	User       *User
	occurredOn time.Time
}

func (e *UserCreatedEvent) EventTopic() EventTopic {
	return "user.created"
}

func (e *UserCreatedEvent) OccurredOn() time.Time {
	return e.occurredOn
}

type RoomCreatedEvent struct {
	Room       *Room
	occurredOn time.Time
}

func (e *RoomCreatedEvent) EventTopic() EventTopic {
	return "room.created"
}

func (e *RoomCreatedEvent) OccurredOn() time.Time {
	return e.occurredOn
}

type RoomJoinedEvent struct {
	UserNumber UserNumber
	RoomNumber RoomNumber
	occurredOn time.Time
}

func (e *RoomJoinedEvent) EventTopic() EventTopic {
	return "room.joined"
}

func (e *RoomJoinedEvent) OccurredOn() time.Time {
	return e.occurredOn
}

type RoomLeftEvent struct {
	UserNumber UserNumber
	RoomNumber RoomNumber
	occurredOn time.Time
}

func (e *RoomLeftEvent) EventTopic() EventTopic {
	return "room.left"
}

func (e *RoomLeftEvent) OccurredOn() time.Time {
	return e.occurredOn
}
