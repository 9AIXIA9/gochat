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
	encryptor.EXPECT().Encrypt(fixedPassword.String()).Return(fixedPasswordEncrypted.String(), nil)
	enc2, err := fixedPassword.Encrypt(encryptor)
	require.NoError(t, err)
	require.Equal(t, fixedPasswordEncrypted, enc2)

	// error path
	encryptorErr := mocks.NewMockEncryptor(ctrl)
	encryptorErr.EXPECT().Encrypt(fixedPassword.String()).Return("", stdErrors.New("encrypt failed"))
	enc3, err := fixedPassword.Encrypt(encryptorErr)
	require.Error(t, err)
	require.Equal(t, domain.PasswordEncrypted(""), enc3)
}

func TestPasswordEncrypted_Compare(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// empty encrypted password always succeeds
	emptyEnc := domain.PasswordEncrypted("")
	err := emptyEnc.Compare(fixedPassword, nil)
	require.NoError(t, err)

	// success path
	comparator := mocks.NewMockComparator(ctrl)
	comparator.EXPECT().Compare(fixedPasswordEncrypted.String(), fixedPassword.String()).Return(nil)
	err = fixedPasswordEncrypted.Compare(fixedPassword, comparator)
	require.NoError(t, err)

	// error path
	comparatorErr := mocks.NewMockComparator(ctrl)
	comparatorErr.EXPECT().Compare(fixedPasswordEncrypted.String(), fixedPassword.String()).Return(stdErrors.New("not match"))
	err = fixedPasswordEncrypted.Compare(fixedPassword, comparatorErr)
	require.ErrorIs(t, err, domain.ErrInvalidPassword)
}
