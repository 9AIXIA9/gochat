package usecase

import (
	"context"
	"gochat/internal/config"
	"gochat/internal/domain"
	"gochat/internal/infra/crypto"
	"gochat/internal/infra/encrypt"
	"gochat/internal/infra/jwt"
	"gochat/internal/utils"
)

type Login struct {
	conf             *config.Token
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
}

func NewLogin(conf *config.Token, useRepo domain.UserRepository, refreshTokenRepo domain.RefreshTokenRepository) domain.LoginUsecase {
	return &Login{
		conf:             conf,
		userRepo:         useRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (uc *Login) Execute(ctx context.Context, req *domain.LoginRequest) (*domain.Response, error) {
	//查询用户信息
	user, err := uc.FindUser(ctx, req.Number)
	if err != nil {
		if utils.IsNotFound(err) {
			return domain.UserNotExistResponse, nil
		}
		return nil, err
	}

	if user == nil {
		return domain.NotJoinedResponse, nil
	}

	// 验证用户名密码
	if err := uc.CheckPwd(ctx, req.Password, user.PwdHash()); err != nil {
		return domain.WrongPasswordResponse, nil
	}

	// 生成 refresh rToken
	rToken, err := uc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	//存储 refresh rToken
	info := &domain.RefreshInfo{Auth: &domain.AuthInfo{UserNumber: req.Number}}

	if err := uc.SaveRefreshToken(ctx, rToken, info); err != nil {
		return nil, err
	}

	//生成 auth token
	aToken, err := uc.GenerateAuthToken(ctx, info.Auth)
	if err != nil {
		return nil, err
	}

	// 返回token
	return domain.NewSuccessResponse(domain.LoginResponse{
		AuthToken:    aToken,
		RefreshToken: rToken,
	}), nil
}

func (uc *Login) FindUser(ctx context.Context, number domain.UserNumber) (*domain.User, error) {
	return uc.userRepo.FindOneByNumber(ctx, number)
}

func (uc *Login) CheckPwd(ctx context.Context, origin, hash string) error {
	return encrypt.Compare(ctx, origin, hash)
}

func (uc *Login) GenerateRefreshToken() (domain.RefreshToken, error) {
	return crypto.GenerateRefreshToken(uc.conf.Refresh.Length)
}

func (uc *Login) SaveRefreshToken(ctx context.Context, token domain.RefreshToken, info *domain.RefreshInfo) error {
	return uc.refreshTokenRepo.Save(ctx, token, info, uc.conf.Refresh.ExpireDuration)
}

func (uc *Login) GenerateAuthToken(ctx context.Context, authInfo *domain.AuthInfo) (domain.AuthToken, error) {
	return jwt.GenerateAuthToken(ctx, uc.conf.Auth.Secret, uc.conf.Auth.ExpireDuration, authInfo)
}
