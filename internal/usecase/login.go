package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gochat/internal/config"
	"gochat/internal/domain"
	"gochat/internal/infra/encrypt"
	"time"
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
		return domain.NewResponseWithoutMsg(domain.CodeUserNotExist), nil
	}

	// 验证用户名密码
	if err := uc.CheckPwd(req.Password, user.PwdHash); err != nil {
		return domain.NewResponseWithoutMsg(domain.CodeWrongPassword), nil
	}

	// 生成token
	token, err := uc.GenerateToken(req.Number)
	if err != nil {
		return nil, err
	}

	// 返回token
	return domain.NewSuccessResponse(gin.H{
		"token":       token,
		"user_number": req.Number,
	}), nil
}

func (uc *Login) QueryUser(number domain.UserNumber) (*domain.User, error) {
	return uc.repo.QueryByNumber(number)
}

func (uc *Login) CheckPwd(origin, hash string) error {
	return encrypt.Compare(origin, hash)
}

func (uc *Login) GenerateToken(userNumber domain.UserNumber) (string, error) {
	mySecret := []byte(uc.conf.Secret)
	dur := time.Duration(uc.conf.ExpireTime) * time.Hour

	c := domain.JwtCustomClaims{
		Auth: &domain.AuthInfo{UserNumber: userNumber},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(dur)),
		},
	}

	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	//使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, err := token.SignedString(mySecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
