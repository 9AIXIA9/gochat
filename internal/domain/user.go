package domain

type UserNumber BaseNumber
type User struct {
	number  UserNumber
	pwdHash string
}

func NewUser(number UserNumber, pwdHash string) *User {
	return &User{
		number:  number,
		pwdHash: pwdHash,
	}
}

func CreateUser(number BaseNumber, pwdHash string) *User {
	return &User{
		number:  UserNumber(number),
		pwdHash: pwdHash,
	}
}

//Getter

func (u *User) Number() UserNumber {
	return u.number
}
func (u *User) PwdHash() string {
	return u.pwdHash
}
