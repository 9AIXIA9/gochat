package usecase

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/utils"
	"time"
)

type Login struct {
	expireDuration time.Duration
	domain.UserFinder
	domain.RefreshTokenSaver
	domain.RefreshTokenGenerator
	domain.AuthTokenGenerator
	domain.Comparator
}

func NewLogin(finder domain.UserFinder,
	saver domain.RefreshTokenSaver,
	refreshTokenGenerator domain.RefreshTokenGenerator,
	authTokenGenerator domain.AuthTokenGenerator,
	comparator domain.Comparator,
	expireDuration time.Duration) domain.LoginUsecase {
	return &Login{
		expireDuration:        expireDuration,
		UserFinder:            finder,
		RefreshTokenSaver:     saver,
		RefreshTokenGenerator: refreshTokenGenerator,
		AuthTokenGenerator:    authTokenGenerator,
		Comparator:            comparator,
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
	if err := uc.Compare(req.Password, user.PwdHash()); err != nil {
		return domain.WrongPasswordResponse, nil
	}

	// 生成 refresh rToken
	rToken, err := uc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	//存储 refresh rToken
	info := &domain.RefreshInfo{Auth: &domain.AuthInfo{UserNumber: req.Number}}

	if err := uc.SaveRefreshToken(ctx, rToken, info, uc.expireDuration); err != nil {
		return nil, err
	}

	//生成 auth token
	aToken, err := uc.GenerateAuthToken(info.Auth)
	if err != nil {
		return nil, err
	}

	// 返回token
	return domain.NewSuccessResponse(domain.LoginResponse{
		AuthToken:    aToken,
		RefreshToken: rToken,
	}), nil
}
