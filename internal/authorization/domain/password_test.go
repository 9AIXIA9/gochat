package domain_test

import (
	"errors"
	"gochat/internal/authorization/domain"
	"gochat/internal/authorization/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedPlainPassword domain.Password          = "plain-pass"
	fixedEncryptedPass domain.PasswordEncrypted = "encrypted-pass"
)

func TestPassword_Encrypt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test successful encryption
	encryptor := mocks.NewMockEncryptor(ctrl)
	encryptor.EXPECT().Encrypt(fixedPlainPassword.String()).Return(fixedEncryptedPass.String(), nil)

	encrypted, err := fixedPlainPassword.Encrypt(encryptor)
	require.NoError(t, err)
	require.Equal(t, fixedEncryptedPass.String(), encrypted.String())

	// Test encrypt error
	errEncryptor := mocks.NewMockEncryptor(ctrl)
	errEncryptor.EXPECT().Encrypt(fixedPlainPassword.String()).Return("", errors.New("encryption error"))

	encrypted2, err := fixedPlainPassword.Encrypt(errEncryptor)
	require.ErrorContains(t, err, "encryption error")
	require.Equal(t, domain.PasswordEncrypted(""), encrypted2)
}

func TestPassword_Validate(t *testing.T) {
	// Test valid password
	validPassword := domain.Password("validPass")
	err := validPassword.Validate()
	require.NoError(t, err)

	// Test too short password
	shortPassword := domain.Password("short")
	err = shortPassword.Validate()
	require.ErrorIs(t, err, myErrors.ErrInvalidLength)

	// Test too long password
	longPassword := domain.Password("this-password-is-was-too-long-to-be-valid")
	err = longPassword.Validate()
	require.ErrorIs(t, err, myErrors.ErrInvalidLength)
}
