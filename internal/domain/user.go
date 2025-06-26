package domain

type UserNumber int64

type User struct {
	Number  UserNumber
	Name    string
	PwdHash string
}

type UserRepository interface {
}

type SignupUsecase interface {
	EncryptPwd(pwd string) string
	GenerateNumber() UserNumber
	CreateUser(user User) error
}

type LoginUsecase interface {
	CheckPwd(number UserNumber, pwd string) error
	GenerateToken() string
}
