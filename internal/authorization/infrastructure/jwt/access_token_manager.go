package jwt

import (
	"errors"
	"fmt"
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var _ domain.AccessTokenParser = (*AccessTokenManager)(nil)
var _ domain.AccessTokenGenerator = (*AccessTokenManager)(nil)

type AccessTokenManager struct {
	secret           string
	validityDuration time.Duration
}

func NewAccessTokenManager(config *AccessTokenConfig) *AccessTokenManager {
	return &AccessTokenManager{
		secret:           config.Secret,
		validityDuration: config.ValidityDuration,
	}
}

type Claims struct {
	UserID kernel.UserID
	jwt.RegisteredClaims
}

func (m *AccessTokenManager) Generate(userID kernel.UserID) (domain.AccessToken, error) {
	mySecret := []byte(m.secret)

	c := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(m.validityDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	//使用指定的secret签名并获得完整的编码后的字符串token
	if tokenStr, err := token.SignedString(mySecret); err != nil {
		return "", fmt.Errorf("get token signed string failed,err:%w", err)
	} else {
		return domain.AccessToken(tokenStr), nil
	}
}

func (m *AccessTokenManager) Parse(accessToken domain.AccessToken) (kernel.UserID, error) {
	//解析 Token
	token, err := jwt.ParseWithClaims(accessToken.String(), &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(m.secret), nil
	})

	if err != nil || !token.Valid {
		return "", ErrInvalidToken
	}

	//验证合理性
	if claims, ok := token.Claims.(*Claims); ok {
		if err := claims.Valid(); err != nil {
			return "", err
		}
		return claims.UserID, nil
	}

	return "", ErrInvalidToken
}

func (c Claims) Valid() error {
	if err := c.RegisteredClaims.Valid(); err != nil {
		return errors.Join(ErrInvalidToken, err)
	}

	if len(c.UserID) == 0 {
		return ErrInvalidToken
	}

	return nil
}
