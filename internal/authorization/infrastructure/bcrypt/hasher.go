package bcrypt

import (
	"errors"
	"fmt"
	"gochat/internal/authorization/application"
	myErrors "gochat/internal/shared/errors"

	"golang.org/x/crypto/bcrypt"
)

var _ application.Comparator = (*Hasher)(nil)
var _ application.Encryptor = (*Hasher)(nil)

type Hasher struct {
	cost int
}

func NewHasher(config *HasherConfig) *Hasher {
	return &Hasher{cost: config.Cost}
}

func (h *Hasher) Compare(hash string, origin string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(origin)); err != nil {
		return errors.Join(myErrors.ErrInvalidCredential, err)
	}
	return nil
}

func (h *Hasher) Encrypt(origin string) (hash string, err error) {
	bits, err := bcrypt.GenerateFromPassword([]byte(origin), h.cost)
	if err != nil {
		return "", fmt.Errorf("encrypt failed,err:%w", err)
	}
	return string(bits), nil
}
