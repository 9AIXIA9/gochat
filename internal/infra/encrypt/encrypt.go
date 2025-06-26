package encrypt

import (
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(origin string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(origin), bcrypt.DefaultCost) //加密处理
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Compare 核对登录密码是否为加密密码
func Compare(origin, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(origin)) //验证（对比）
}
