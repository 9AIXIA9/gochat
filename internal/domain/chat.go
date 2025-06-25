package domain

import "time"

type Chat struct {
	UserNumber UserNumber
	RoomNumber RoomNumber
	Content    string
	SendTime   time.Time
}

type SendMsg interface {
	RecordChat(chat Chat) error
	Broadcast(chat Chat) error
}

type LoadChats interface {
	LoadChats(number RoomNumber) ([]Chat, error)
}
