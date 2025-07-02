package domain

type UserNumber int64
type User struct {
	Number  UserNumber
	Name    string
	PwdHash string
	Events  []Event
}
