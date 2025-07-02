package domain

import "time"

type Chat struct {
	UserNumber UserNumber
	RoomNumber RoomNumber
	Content    string
	SendTime   time.Time
	Events     []Event
}
