package domain

import "time"

type Chat struct {
	userNumber UserNumber
	roomNumber RoomNumber
	content    string
	sendTime   time.Time
	events     []Event
}

func NewChat(userNumber UserNumber, roomNumber RoomNumber, content string, sendTime time.Time) *Chat {
	return &Chat{
		userNumber: userNumber,
		roomNumber: roomNumber,
		content:    content,
		sendTime:   sendTime,
		events:     make([]Event, 0),
	}
}

// GetEvents - 获取领域事件
func (c *Chat) GetEvents() []Event {
	return c.events
}

// ClearEvents - 清空领域事件（通常在事件发布后调用）
func (c *Chat) ClearEvents() {
	c.events = make([]Event, 0)
}

// addEvent - 添加领域事件
func (c *Chat) addEvent(event Event) {
	c.events = append(c.events, event)
}

//Getter

func (c *Chat) UserNumber() UserNumber {
	return c.userNumber
}

func (c *Chat) RoomNumber() RoomNumber {
	return c.roomNumber
}

func (c *Chat) Content() string {
	return c.content
}

func (c *Chat) SendTime() time.Time {
	return c.sendTime
}
