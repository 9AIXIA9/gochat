package usecase

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/domain"
	"gochat/internal/infra/encrypt"
	"gochat/internal/infra/snowflake"
)

type Signup struct {
	repo domain.UserRepository
}

func NewSignup(repo domain.UserRepository) domain.SignupUsecase {
	return &Signup{repo: repo}
}

func (uc *Signup) Logic(req *domain.SignupRequest) (*domain.Response, error) {
	// 加密密码
	hashedPassword, err := uc.EncryptPwd(req.Password)
	if err != nil {
		return nil, err
	}

	// 生成用户号码
	userNumber := uc.GenerateNumber()

	// 创建用户
	user := &domain.User{
		Number:  userNumber,
		Name:    req.Name,
		PwdHash: hashedPassword,
	}

	if exist, err := uc.CreateUser(user); err != nil {
		return nil, err
	} else if exist {
		return domain.NewResponseWithoutMsg(domain.CodeUserExist), nil
	}

	// 返回用户信息
	return domain.NewSuccessResponse(gin.H{
		"username":    req.Name,
		"user_number": userNumber,
	}), nil
}

func (uc *Signup) EncryptPwd(pwd string) (string, error) {
	return encrypt.Encrypt(pwd)
}

func (uc *Signup) GenerateNumber() domain.UserNumber {
	return snowflake.GenerateUserNumber()
}

func (uc *Signup) CreateUser(user *domain.User) (bool, error) {
	return uc.repo.Create(user)
}
