package domain_test

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/authorization/domain/mocks"
	eventMock "gochat/internal/shared/event/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoadUser(t *testing.T) {
	start := time.Now().UTC().
		UTC()
	user := domain.LoadUser(
		fixedUserID,
		fixedEmail,
		fixedUserNumber,
		fixedPasswordEncrypted,
		time.Now().UTC(),
	)
	require.NotNil(t, user)
	assert.Equal(t, fixedUserID, user.ID())
	assert.Equal(t, fixedEmail, user.Email())
	assert.Equal(t, fixedUserNumber, user.Number())
	assert.Equal(t, fixedPasswordEncrypted, user.PasswordEncrypted())
	assert.WithinDuration(t, start, user.SignedUpAt(), time.Second)
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserIDGenerator := mocks.NewMockUserIDGenerator(ctrl)
	mockNumberGenerator := mocks.NewMockUserNumberGenerator(ctrl)
	mockEventIDGenerator := eventMock.NewMockIDGenerator(ctrl)

	gomock.InOrder(
		mockUserIDGenerator.EXPECT().Generate().Return(fixedUserID).Times(1),
		mockNumberGenerator.EXPECT().Generate().Return(fixedUserNumber).Times(1),
		mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1),
	)
	user, err := domain.CreateUser(
		fixedEmail,
		fixedPasswordEncrypted,
		mockUserIDGenerator,
		mockNumberGenerator,
		mockEventIDGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, fixedUserID, user.ID())
	assert.Equal(t, fixedEmail, user.Email())
	assert.Equal(t, fixedUserNumber, user.Number())
	assert.Equal(t, fixedPasswordEncrypted, user.PasswordEncrypted())
	assert.WithinDuration(t, time.Now().UTC().
		UTC(), user.SignedUpAt(), timeTolerance)

	evs := user.GetEvents()
	require.Len(t, evs, 1)

	evCreated := evs[0]
	assert.Equal(t, domain.TopicUserCreated, evCreated.Topic())
	assert.Equal(t, fixedUserID.String(), evCreated.AggregateID().String())
}
