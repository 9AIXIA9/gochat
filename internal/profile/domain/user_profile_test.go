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

const (
	maxNameLen    = 32
	maxAddressLen = 100
	maxSignLen    = 100
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

func TestUserProfile_Updates(t *testing.T) {
	profile := domain.LoadUserProfile(
		fixedUserID,
		"",
		kernel.UnknownGender,
		"",
		"",
		"",
		"",
		time.Now().UTC(),
	)

	err := profile.UpdateName(fixedName)
	require.NoError(t, err)
	assert.Equal(t, fixedName, profile.Name())

	profile.UpdateGender(fixedGender)
	assert.Equal(t, fixedGender, profile.Gender())

	profile.UpdateEmail(fixedEmail)
	assert.Equal(t, fixedEmail, profile.Email())

	profile.UpdatePhoneNumber(fixedPhoneNumber)
	assert.Equal(t, fixedPhoneNumber, profile.PhoneNumber())

	err = profile.UpdateAddress(fixedAddress)
	require.NoError(t, err)
	assert.Equal(t, fixedAddress, profile.Address())

	err = profile.UpdateSign(fixedSign)
	require.NoError(t, err)
	assert.Equal(t, fixedSign, profile.Sign())

	// Test exceeding max lengths
	longName := ""
	for i := 0; i < maxNameLen+1; i++ {
		longName += "a"
	}
	err = profile.UpdateName(longName)
	require.Error(t, err)

	longAddress := ""
	for i := 0; i < maxAddressLen+1; i++ {
		longAddress += "a"
	}
	err = profile.UpdateAddress(kernel.Address(longAddress))
	require.Error(t, err)

	longSign := ""
	for i := 0; i < maxSignLen+1; i++ {
		longSign += "a"
	}
	err = profile.UpdateSign(longSign)
	require.Error(t, err)
}
