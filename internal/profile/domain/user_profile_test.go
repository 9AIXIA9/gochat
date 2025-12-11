package domain_test

import (
	"gochat/internal/profile/domain"
	"gochat/internal/profile/domain/mocks"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoadUserProfile(t *testing.T) {
	profile := domain.LoadUserProfile(
		fixedProfileID,
		fixedName,
		fixedGender,
		fixedEmail,
		fixedPhoneNumber,
		&kernel.Address{
			Country:  fixedCountry,
			Province: fixedProvince,
			City:     fixedCity,
			District: fixedDistrict,
			Street:   fixedStreet,
		},
		fixedSign,
		time.Now().UTC(),
	)

	require.NotNil(t, profile)
	assert.Equal(t, fixedProfileID, profile.ID())
	assert.Equal(t, fixedName, profile.Name())
	assert.Equal(t, fixedGender, profile.Gender())
	assert.Equal(t, fixedEmail, profile.Email())
	assert.Equal(t, fixedPhoneNumber, profile.PhoneNumber())
	assert.NotNil(t, profile.Address())
	assert.Equal(t, fixedCountry, profile.Address().Country)
	assert.Equal(t, fixedProvince, profile.Address().Province)
	assert.Equal(t, fixedCity, profile.Address().City)
	assert.Equal(t, fixedDistrict, profile.Address().District)
	assert.Equal(t, fixedStreet, profile.Address().Street)
	assert.Equal(t, fixedSign, profile.Sign())
}

func TestCreateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := mocks.NewMockProfileIDGenerator(ctrl)
	mockIDGenerator.EXPECT().
		Generate().
		Return(fixedProfileID).
		Times(1)

	signedUpAt := time.Now().UTC()
	profile := domain.CreateUserProfile(
		fixedEmail,
		signedUpAt,
		mockIDGenerator,
	)

	require.NotNil(t, profile)
	assert.Equal(t, fixedProfileID, profile.ID())
	assert.Equal(t, fixedEmail, profile.Email())
	assert.Equal(t, signedUpAt, profile.SignedUpAt())
	assert.Equal(t, domain.ProfileID("profile-789"), profile.ID())
	assert.Equal(t, kernel.UnknownGender, profile.Gender())
}
