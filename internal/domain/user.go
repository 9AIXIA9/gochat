package domain

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
	GenerateToken(authInfo *AuthInfo) (string, error)
}
