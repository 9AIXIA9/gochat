package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"gochat/internal/authorization/domain"
)

var _ domain.RandomStringGenerator = (*RandomStringGenerator)(nil)

type RandomStringGenerator struct {
	length int
}

func NewRefreshTokenGenerator(config *RefreshTokenConfig) *RandomStringGenerator {
	return &RandomStringGenerator{length: config.Length}
}

// Generate 生成一个安全的随机令牌
func (g *RandomStringGenerator) Generate() (string, error) {
	b := make([]byte, g.length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
