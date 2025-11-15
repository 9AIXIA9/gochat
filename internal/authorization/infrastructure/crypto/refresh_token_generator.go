package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
)

var _ application.RefreshTokenGenerator = (*RefreshTokenGenerator)(nil)

type RefreshTokenGenerator struct {
	length int
}

func NewRefreshTokenGenerator(config *RefreshTokenConfig) *RefreshTokenGenerator {
	return &RefreshTokenGenerator{length: config.Length}
}

// Generate 生成一个安全的随机令牌
func (g *RefreshTokenGenerator) Generate() (domain.RefreshToken, error) {
	b := make([]byte, g.length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return domain.RefreshToken(base64.RawURLEncoding.EncodeToString(b)), nil
}
