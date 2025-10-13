package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"gochat/internal/domain"
	"gochat/internal/types"
)

var _ domain.RefreshTokenGenerator = (*RefreshTokenGenerator)(nil)

type RefreshTokenGenerator struct {
	length int
}

func NewRefreshTokenGenerator(length int) (*RefreshTokenGenerator, error) {
	if length <= 0 {
		return nil, types.ErrLengthLessThanZero
	}

	return &RefreshTokenGenerator{length: length}, nil
}

// GenerateRefreshToken 生成一个安全的随机令牌
func (g *RefreshTokenGenerator) GenerateRefreshToken() (domain.RefreshToken, error) {
	b := make([]byte, g.length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return domain.RefreshToken(base64.RawURLEncoding.EncodeToString(b)), nil
}
