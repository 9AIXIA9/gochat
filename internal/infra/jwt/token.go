package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"gochat/internal/config"
	"gochat/internal/domain"
	"time"
)

type Claims struct {
	UserNumber domain.UserNumber
	jwt.RegisteredClaims
}

func GenerateToken(conf *config.JWT, info *domain.AuthInfo) (string, error) {
	mySecret := []byte(conf.Secret)
	dur := time.Duration(conf.ExpireTime) * time.Hour

	c := Claims{
		UserNumber: info.UserNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(dur)),
		},
	}

	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	//使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

func ParseToken(secret, tokenStr string) (*domain.AuthInfo, error) {
	//解析 Token
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	//验证合理性
	if claims, ok := token.Claims.(*Claims); ok {
		if err := claims.Valid(); err != nil {
			return nil, err
		}
		return &domain.AuthInfo{UserNumber: claims.UserNumber}, nil
	}

	return nil, domain.ErrInvalidTokenClaims
}

func (c Claims) Valid() error {
	if err := c.RegisteredClaims.Valid(); err != nil {
		return err
	}

	if c.UserNumber == 0 {
		return domain.ErrInvalidTokenClaims
	}

	return nil
}
