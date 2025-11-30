package domain_test

import (
	stdErrors "errors"
	"testing"

	"gochat/internal/roomship/domain"
	"gochat/internal/roomship/domain/mocks"
	sharedErrors "gochat/internal/shared/errors"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedRawPassword       domain.Password = "abc123"
	fixedEncryptedPassword                 = "encrypted-abc123"
)

func TestPassword_Validate(t *testing.T) {
	// exactly max length (20) OK
	p := domain.Password("12345678901234567890")
	require.NoError(t, p.Validate())

	// exceed max length (>20)
	p2 := domain.Password("123456789012345678901")
	require.ErrorIs(t, p2.Validate(), sharedErrors.ErrInvalidLength)
}

func TestPassword_Encrypt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// empty password returns empty encrypted without calling encryptor (nil acceptable)
	empty := domain.Password("")
	enc, err := empty.Encrypt(nil)
	require.NoError(t, err)
	require.Equal(t, domain.PasswordEncrypted(""), enc)

	// success path
	encryptor := mocks.NewMockEncryptor(ctrl)
	encryptor.EXPECT().Encrypt(fixedRawPassword.String()).Return(fixedEncryptedPassword, nil)
	enc2, err := fixedRawPassword.Encrypt(encryptor)
	require.NoError(t, err)
	require.Equal(t, domain.PasswordEncrypted(fixedEncryptedPassword), enc2)

	// error path
	encryptorErr := mocks.NewMockEncryptor(ctrl)
	encryptorErr.EXPECT().Encrypt(fixedRawPassword.String()).Return("", stdErrors.New("encrypt failed"))
	enc3, err := fixedRawPassword.Encrypt(encryptorErr)
	require.Error(t, err)
	require.Equal(t, domain.PasswordEncrypted(""), enc3)
}
