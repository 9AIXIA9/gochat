package domain

import "time"

type UserNumber int64

type RoomNumber int64

type User struct {
	Number  UserNumber
	Name    string
	PwdHash string
}

type Room struct {
	Name         string
	Number       RoomNumber
	SecretHash   string
	Description  string
	CurrentUsers int
	MaxUsers     int
	Owner        UserNumber
}

type Chat struct {
	UserNumber UserNumber
	RoomNumber RoomNumber
	Content    string
	SendTime   time.Time
}
