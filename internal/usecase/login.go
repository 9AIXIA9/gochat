package usecase

import (
	"gochat/internal/domain"
)

type Login struct{}

func NewLogin() domain.LoginUsecase {
	return &Login{}
}

func (uc *Login) CheckPwd(number domain.UserNumber, pwd string) error {
	//TODO implement me
	panic("implement me")
}

func (uc *Login) GenerateToken() string {
	//TODO implement me
	panic("implement me")
}
