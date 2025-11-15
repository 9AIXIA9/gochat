package event

type Manager struct {
	events []Event
}

func NewEventManager() *Manager {
	return &Manager{
		events: make([]Event, 0),
	}
}

// GetEvents 获取并清空待发布的事件
func (e *Manager) GetEvents() []Event {
	events := e.events
	e.events = nil
	return events
}

// RecordEvent 记录事件
func (e *Manager) RecordEvent(event Event) {
	if e.events == nil {
		e.events = make([]Event, 0, 1) //提前分配内存
	}
	e.events = append(e.events, event)
}
