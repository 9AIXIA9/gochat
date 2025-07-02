package domain

type RoomNumber int64
type Room struct {
	Name         string
	Number       RoomNumber
	SecretHash   string
	Description  string
	CurrentUsers int
	MaxUsers     int
	Owner        UserNumber
	Events       []Event
}
