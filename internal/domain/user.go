package domain

type UserNumber int64
type User struct {
	number  UserNumber
	name    string
	pwdHash string
}

func NewUser(number UserNumber, name string, pwdHash string) *User {
	return &User{
		number:  number,
		name:    name,
		pwdHash: pwdHash,
	}
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
