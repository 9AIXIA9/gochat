package usecase

import (
	"context"
	"gochat/internal/config"
	"gochat/internal/domain"
	"gochat/internal/infra/crypto"
	"gochat/internal/infra/jwt"
)

type RefreshToken struct {
	conf *config.Token
	repo domain.RefreshTokenRepository
}

func NewRefreshToken(conf *config.Token, repo domain.RefreshTokenRepository) domain.RefreshTokenUsecase {
	return &RefreshToken{
		conf: conf,
		repo: repo,
	}
}

func (uc *RefreshToken) Execute(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.Response, error) {
	//查询 refresh token是否存在
	info, err := uc.QueryRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return domain.UnauthorizedResponse, nil
	}

	//生成 refresh token
	rToken, err := uc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	//存储 refresh token
	if err = uc.SaveRefreshToken(ctx, rToken, info); err != nil {
		return nil, err
	}

	//生成 auth token
	aToken, err := uc.GenerateAuthToken(ctx, info.Auth)
	if err != nil {
		return nil, err
	}
	return domain.NewSuccessResponse(domain.RefreshTokenResponse{
		AuthToken:    aToken,
		RefreshToken: rToken,
	}), nil
}

func (uc *RefreshToken) QueryRefreshToken(ctx context.Context, token domain.RefreshToken) (*domain.RefreshInfo, error) {
	return uc.repo.FindByToken(ctx, token)
}

func (uc *RefreshToken) SaveRefreshToken(ctx context.Context, token domain.RefreshToken, info *domain.RefreshInfo) error {
	return uc.repo.Save(ctx, token, info, uc.conf.Refresh.ExpireDuration)
}

func (uc *RefreshToken) GenerateRefreshToken() (domain.RefreshToken, error) {
	return crypto.GenerateRefreshToken(uc.conf.Refresh.Length)
}

func (uc *RefreshToken) GenerateAuthToken(ctx context.Context, authInfo *domain.AuthInfo) (domain.AuthToken, error) {
	return jwt.GenerateAuthToken(ctx, uc.conf.Auth.Secret, uc.conf.Auth.ExpireDuration, authInfo)
}
