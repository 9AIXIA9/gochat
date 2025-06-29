package usecase

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/config"
	"gochat/internal/domain"
	"gochat/internal/infra/encrypt"
	"gochat/internal/infra/jwt"
)

type Login struct {
	repo domain.UserRepository
	conf *config.JWT
}

func NewLogin(conf *config.JWT, repo domain.UserRepository) domain.LoginUsecase {
	return &Login{repo: repo, conf: conf}
}

func (uc *Login) Logic(req *domain.LoginRequest) (*domain.Response, error) {
	//查询用户信息
	user, err := uc.QueryUser(req.Number)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return domain.NewResponseWithDefaultMsg(domain.CodeUserNotExist), nil
	}

	// 验证用户名密码
	if err := uc.CheckPwd(req.Password, user.PwdHash); err != nil {
		return domain.NewResponseWithDefaultMsg(domain.CodeWrongPassword), nil
	}

	// 生成token
	token, err := uc.GenerateToken(&domain.AuthInfo{UserNumber: req.Number})
	if err != nil {
		return nil, err
	}

	// 返回token
	return domain.NewSuccessResponse(gin.H{
		"token": token,
	}), nil
}

func (uc *Login) QueryUser(number domain.UserNumber) (*domain.User, error) {
	return uc.repo.QueryByNumber(number)
}

func (uc *Login) CheckPwd(origin, hash string) error {
	return encrypt.Compare(origin, hash)
}

func (uc *Login) GenerateToken(authInfo *domain.AuthInfo) (string, error) {
	return jwt.GenerateToken(uc.conf, authInfo)
}
