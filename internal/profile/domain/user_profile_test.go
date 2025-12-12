package domain_test

import (
	"gochat/internal/profile/domain"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoadUserProfile(t *testing.T) {
	profile := domain.LoadUserProfile(
		fixedUserID,
		fixedName,
		fixedGender,
		fixedEmail,
		fixedPhoneNumber,
		fixedAddress,
		fixedSign,
		time.Now().UTC(),
	)

	require.NotNil(t, profile)
	assert.Equal(t, fixedUserID, profile.ID())
	assert.Equal(t, fixedName, profile.Name())
	assert.Equal(t, fixedGender, profile.Gender())
	assert.Equal(t, fixedEmail, profile.Email())
	assert.Equal(t, fixedPhoneNumber, profile.PhoneNumber())
	assert.Equal(t, fixedAddress, profile.Address())
	assert.Equal(t, fixedSign, profile.Sign())
}

func TestCreateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	signedUpAt := time.Now().UTC()
	profile := domain.CreateUserProfile(
		fixedUserID,
		fixedEmail,
		signedUpAt,
	)

	require.NotNil(t, profile)
	assert.Equal(t, fixedUserID, profile.ID())
	assert.Equal(t, fixedEmail, profile.Email())
	assert.Equal(t, signedUpAt, profile.SignedUpAt())
	assert.Equal(t, fixedUserID, profile.ID())
	assert.Equal(t, kernel.UnknownGender, profile.Gender())
}
