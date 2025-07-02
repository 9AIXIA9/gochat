package usecase

import (
	"context"
	"gochat/internal/config"
	"gochat/internal/domain"
	"gochat/internal/infra/encrypt"
	"gochat/internal/infra/jwt"
	"gochat/internal/utils"
)

type Login struct {
	repo domain.UserRepository
	conf *config.JWT
}

func NewLogin(conf *config.JWT, repo domain.UserRepository) domain.LoginUsecase {
	return &Login{repo: repo, conf: conf}
}

func (uc *Login) Logic(ctx context.Context, req *domain.LoginRequest) (*domain.Message, error) {
	//查询用户信息
	user, err := uc.FindUser(ctx, req.Number)
	if err != nil {
		if utils.IsNotFound(err) {
			return domain.UserNotExistMessage, nil
		}
		return nil, err
	}

	if user == nil {
		return domain.NotJoinedMessage, nil
	}

	// 验证用户名密码
	if err := uc.CheckPwd(ctx, req.Password, user.PwdHash()); err != nil {
		return domain.WrongPasswordMessage, nil
	}

	// 生成token
	token, err := uc.GenerateToken(ctx, &domain.AuthInfo{UserNumber: req.Number})
	if err != nil {
		return nil, err
	}

	// 返回token
	return domain.NewSuccessMessage(domain.LoginResponse{Token: token}), nil
}

func (uc *Login) FindUser(ctx context.Context, number domain.UserNumber) (*domain.User, error) {
	return uc.repo.FindOneByNumber(ctx, number)
}

func (uc *Login) CheckPwd(ctx context.Context, origin, hash string) error {
	return encrypt.Compare(ctx, origin, hash)
}

func (uc *Login) GenerateToken(ctx context.Context, authInfo *domain.AuthInfo) (string, error) {
	return jwt.GenerateToken(ctx, uc.conf, authInfo)
}
