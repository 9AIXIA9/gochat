package model

type DeadEvent struct {
	*Event
	Reason string
}

func NewDeadEvent(event *Event, reason error) *DeadEvent {
	return &DeadEvent{Event: event, Reason: reason.Error()}
}

func (e *DeadEvent) TableName() string {
	return "gochat.dead_events"
}
