package domain

import "time"

type UserNumber int64
type User struct {
	number  UserNumber
	name    string
	pwdHash string
	events  []Event
}

func NewUser(number UserNumber, name string, pwdHash string) *User {
	return &User{
		number:  number,
		name:    name,
		pwdHash: pwdHash,
		events:  make([]Event, 0),
	}
}

func CreateUser(number UserNumber, name string, pwdHash string) *User {
	u := NewUser(number, name, pwdHash)
	u.addEvent(&UserCreatedEvent{
		User:       u,
		occurredOn: time.Now(),
	})
	return u
}

// GetEvents - 获取领域事件
func (u *User) GetEvents() []Event {
	return u.events
}

// ClearEvents - 清空领域事件（通常在事件发布后调用）
func (u *User) ClearEvents() {
	u.events = make([]Event, 0)
}

// addEvent - 添加领域事件
func (u *User) addEvent(event Event) {
	u.events = append(u.events, event)
}

//Getter

func (u *User) Number() UserNumber {
	return u.number
}
func (u *User) Name() string {
	return u.name
}
func (u *User) PwdHash() string {
	return u.pwdHash
}
