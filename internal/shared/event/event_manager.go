package event

//TODO 可以做对象池

type Manager struct {
	events []Event
}

func NewEventManager() *Manager {
	return &Manager{
		events: nil,
	}
}

// GetEvents 获取并清空待发布的事件
func (e *Manager) GetEvents() []Event {
	events := e.events
	e.events = nil
	return events
}

// RecordEvent 记录事件
func (e *Manager) RecordEvent(event SpecificEvent) {
	if e.events == nil {
		e.events = make([]Event, 0, 1)
	}
	e.events = append(e.events, event)
}
