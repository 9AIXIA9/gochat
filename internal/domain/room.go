package domain

import (
	"time"
)

type RoomNumber int64
type Room struct {
	name         string
	number       RoomNumber
	secretHash   string
	description  string
	currentUsers int
	maxUsers     int
	owner        UserNumber
	events       []Event
}

func NewRoom(number RoomNumber, name string, secretHash string, description string, currentUsers, maxUsers int, owner UserNumber) *Room {
	return &Room{
		number:       number,
		name:         name,
		secretHash:   secretHash,
		description:  description,
		currentUsers: currentUsers,
		maxUsers:     maxUsers,
		owner:        owner,
		events:       make([]Event, 0),
	}
}

func CreateRoom(number RoomNumber, name string, secretHash string, description string, maxUsers int, owner UserNumber) *Room {
	r := NewRoom(number, name, secretHash, description, 1, maxUsers, owner)
	r.addEvent(&RoomCreatedEvent{
		Room:       r,
		occurredOn: time.Now(),
	})
	return r
}

// GetEvents - 获取领域事件
func (r *Room) GetEvents() []Event {
	return r.events
}

// ClearEvents - 清空领域事件（通常在事件发布后调用）
func (r *Room) ClearEvents() {
	r.events = make([]Event, 0)
}

// addEvent - 添加领域事件
func (r *Room) addEvent(event Event) {
	r.events = append(r.events, event)
}

//Getter

func (r *Room) Number() RoomNumber {
	return r.number
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) SecretHash() string {
	return r.secretHash
}

func (r *Room) Description() string {
	return r.description
}

func (r *Room) CurrentUsers() int {
	return r.currentUsers
}

func (r *Room) MaxUsers() int {
	return r.maxUsers
}

func (r *Room) Owner() UserNumber {
	return r.owner
}
