package domain

type UserNumber int64

type User struct {
	Number  UserNumber
	Name    string
	PwdHash string
}

type UserRepository interface {
	Create(user *User) (bool, error)
	QueryByNumber(number UserNumber) (*User, error)
}

type SignupUsecase interface {
	Logic(req *SignupRequest) (*Response, error)
	EncryptPwd(pwd string) (string, error)
	GenerateNumber() UserNumber
	CreateUser(user *User) (bool, error)
}

type LoginUsecase interface {
	Logic(req *LoginRequest) (*Response, error)
	QueryUser(number UserNumber) (*User, error)
	CheckPwd(origin, hash string) error
	GenerateToken(userNumber UserNumber) (string, error)
}
