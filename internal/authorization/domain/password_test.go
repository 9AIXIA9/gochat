package domain_test

import (
	"errors"
	"gochat/internal/authorization/domain"
	"gochat/internal/authorization/domain/mocks"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPassword_Encrypt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test successful encryption
	encryptor := mocks.NewMockEncryptor(ctrl)
	encryptor.EXPECT().Encrypt(fixedPassword.String()).Return(fixedPasswordEncrypted.String(), nil)

	encrypted, err := fixedPassword.Encrypt(encryptor)
	require.NoError(t, err)
	require.Equal(t, fixedPasswordEncrypted.String(), encrypted.String())

	// Test encrypt error
	errEncryptor := mocks.NewMockEncryptor(ctrl)
	errEncryptor.EXPECT().Encrypt(fixedPassword.String()).Return("", errors.New("encryption error"))

	encrypted2, err := fixedPassword.Encrypt(errEncryptor)
	require.ErrorContains(t, err, "encryption error")
	require.Equal(t, domain.PasswordEncrypted(""), encrypted2)
}

func TestPasswordEncrypted_Compare(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test successful comparison
	comparator := mocks.NewMockComparator(ctrl)
	comparator.EXPECT().Compare(fixedPasswordEncrypted.String(), fixedPassword.String()).Return(nil)

	err := fixedPasswordEncrypted.Compare(fixedPassword, comparator)
	require.NoError(t, err)

	// Test non-matching password
	comparatorNonMatch := mocks.NewMockComparator(ctrl)
	comparatorNonMatch.EXPECT().Compare(fixedPasswordEncrypted.String(), fixedPassword.String()).Return(errors.New("test"))

	err = fixedPasswordEncrypted.Compare(fixedPassword, comparatorNonMatch)
	require.Error(t, err)
}

func TestPassword_Validate(t *testing.T) {
	// Test valid password
	validPassword := domain.Password("validPass")
	err := validPassword.Validate()
	require.NoError(t, err)

	// Test too short password
	shortPassword := domain.Password("short")
	err = shortPassword.Validate()
	require.Error(t, err)

	// Test too long password
	longPassword := domain.Password("this-password-is-was-too-long-to-be-valid")
	err = longPassword.Validate()
	require.Error(t, err)
}
