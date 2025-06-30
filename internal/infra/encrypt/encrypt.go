package encrypt

import (
	"context"
	"gochat/internal/utils/timeout"
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(ctx context.Context, origin string) (string, error) {
	return timeout.ConvertAndExecuteWithResponse(ctx, func() (string, error) {
		hash, err := bcrypt.GenerateFromPassword([]byte(origin), bcrypt.DefaultCost) //加密处理
		if err != nil {
			return "", err
		}
		return string(hash), nil
	})
}

// Compare 核对登录密码是否为加密密码
func Compare(ctx context.Context, origin, hash string) error {
	return timeout.ConvertAndExecute(ctx, func() error {
		return bcrypt.CompareHashAndPassword([]byte(hash), []byte(origin))
	})
}
