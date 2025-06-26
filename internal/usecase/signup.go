package usecase

import (
	"gochat/internal/domain"
)

type Signup struct {
}

func NewSignup() domain.SignupUsecase {
	return &Signup{}
}

func (uc *Signup) EncryptPwd(pwd string) string {
	//TODO implement me
	panic("implement me")
}

func (uc *Signup) GenerateNumber() domain.UserNumber {
	//TODO implement me
	panic("implement me")
}

func (uc *Signup) CreateUser(user domain.User) error {
	//TODO implement me
	panic("implement me")
}
