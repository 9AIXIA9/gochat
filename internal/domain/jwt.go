package domain

import "github.com/golang-jwt/jwt/v4"

const AuthKey = "auth"

// AuthContext 认证上下文接口
type AuthContext interface {
	SetUserNumber(number UserNumber)
}

// AuthInfo 身份验证信息
type AuthInfo struct {
	UserNumber UserNumber `json:"-"`
}

func (r *AuthInfo) SetUserNumber(number UserNumber) {
	r.UserNumber = number
}

type JwtCustomClaims struct {
	Auth *AuthInfo
	jwt.RegisteredClaims
}
