package domain

type RoomNumber uint64

type Room struct {
	Name       string
	Number     RoomNumber
	SecretHash string
}

type CreateRoom interface {
	EncryptSecret(secret string) string
	GenerateNumber() RoomNumber
	CreateRoom(room Room) error
}

type JoinRoom interface {
	CheckSecret(number RoomNumber, secret string) error
	EstablishLongConnection() error
}

type ExitRoom interface {
	CloseLongConnection() error
}
