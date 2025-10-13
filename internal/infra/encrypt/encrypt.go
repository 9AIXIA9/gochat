package encrypt

import (
	"gochat/internal/config"
	"gochat/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var _ domain.Encryptor = (*BcryptEncryptor)(nil)
var _ domain.Comparator = (*BcryptEncryptor)(nil)

type BcryptEncryptor struct {
	cost int
}

func NewBcryptEncryptor(conf *config.Encryptor) *BcryptEncryptor {
	if conf.Cost <= 0 {
		conf.Cost = bcrypt.DefaultCost
	}
	return &BcryptEncryptor{cost: conf.Cost}
}

func (b *BcryptEncryptor) Encrypt(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), b.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (b *BcryptEncryptor) Compare(origin, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(origin))
}
