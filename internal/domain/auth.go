package domain

import "context"

const AuthInfoKey = "auth"

// AuthContext 认证上下文接口
type AuthContext interface {
	SetInfo(info *AuthInfo)
}

// AuthInfo 身份验证信息
type AuthInfo struct {
	UserNumber UserNumber
}

func (r *AuthInfo) SetInfo(info *AuthInfo) {
	r.UserNumber = info.UserNumber
}

//todo Refresh token

type AuthUsecase interface {
	ParseToken(ctx context.Context, tokenStr string) (*AuthInfo, error)
}
