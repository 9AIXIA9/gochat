package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"gochat/internal/domain"
	"gochat/internal/types"
)

// GenerateRefreshToken 生成一个安全的随机令牌
func GenerateRefreshToken(length int) (domain.RefreshToken, error) {
	if length <= 0 {
		return "", types.ErrLengthLessThanZero
	}
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return domain.RefreshToken(base64.RawURLEncoding.EncodeToString(b)), nil
}
