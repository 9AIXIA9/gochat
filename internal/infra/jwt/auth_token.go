package jwt

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"gochat/internal/domain"
	"gochat/internal/types"
	"gochat/internal/utils/timeout"
	"time"
)

type Claims struct {
	UserNumber domain.UserNumber
	jwt.RegisteredClaims
}

func GenerateAuthToken(ctx context.Context, secret string, expireDuration time.Duration, info *domain.AuthInfo) (domain.AuthToken, error) {
	return timeout.ConvertAndExecuteWithResponse(ctx, func() (domain.AuthToken, error) {
		mySecret := []byte(secret)

		c := Claims{
			UserNumber: info.UserNumber,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			},
		}

		//使用指定的签名方法创建签名对象
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

		//使用指定的secret签名并获得完整的编码后的字符串token
		if tokenStr, err := token.SignedString(mySecret); err != nil {
			return "", nil
		} else {
			return domain.AuthToken(tokenStr), nil
		}
	})
}

func ParseAuthToken(ctx context.Context, secret string, authToken domain.AuthToken) (*domain.AuthInfo, error) {
	return timeout.ConvertAndExecuteWithResponse(ctx, func() (*domain.AuthInfo, error) {
		//解析 Token
		token, err := jwt.ParseWithClaims(string(authToken), &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return nil, types.ErrInvalidToken
		}

		//验证合理性
		if claims, ok := token.Claims.(*Claims); ok {
			if err := claims.Valid(); err != nil {
				return nil, err
			}
			return &domain.AuthInfo{UserNumber: claims.UserNumber}, nil
		}

		return nil, types.ErrInvalidTokenClaims
	})
}

func (c Claims) Valid() error {
	if err := c.RegisteredClaims.Valid(); err != nil {
		return err
	}

	if c.UserNumber == 0 {
		return types.ErrInvalidTokenClaims
	}

	return nil
}
