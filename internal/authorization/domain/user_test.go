package domain_test

import (
	"errors"
	"testing"
	"time"

	"gochat/internal/authorization/domain"
	domainMocks "gochat/internal/authorization/domain/mocks"
	"gochat/internal/shared/event"
	eventMock "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// helper: time tolerance
const timeTolerance = 150 * time.Millisecond

// fixed values used in tests
const (
	fixedUserID    kernel.UserID            = "user-123"
	fixedUserNum   kernel.UserNumber        = "10001"
	fixedEventID   event.ID                 = "ev-999"
	fixedEmail     kernel.Email             = "test@example.com"
	fixedEncrypted domain.PasswordEncrypted = "encrypted-pass"
)

// CreateUser success scenario
func TestUser_CreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	eventIDGenerator.EXPECT().Generate().Return(fixedEventID)

	userIDGenerator := domainMocks.NewMockUserIDGenerator(ctrl)
	userIDGenerator.EXPECT().Generate().Return(fixedUserID)

	userNumberGenerator := domainMocks.NewMockUserNumberGenerator(ctrl)
	userNumberGenerator.EXPECT().Generate().Return(fixedUserNum)

	start := time.Now().UTC()

	user, err := domain.CreateUser(
		fixedEmail,
		fixedEncrypted,
		userIDGenerator,
		userNumberGenerator,
		eventIDGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, user)

	assert.Equal(t, fixedUserID, user.ID())
	assert.Equal(t, fixedUserNum, user.Number())
	assert.Equal(t, fixedEmail, user.Email())
	assert.Equal(t, fixedEncrypted, user.PasswordEncrypted())
	assert.WithinDuration(t, time.Now().UTC(), user.SignedUpAt(), timeTolerance)
	assert.True(t, !user.SignedUpAt().Before(start))

	events := user.GetEvents()
	assert.Len(t, events, 1)
	assert.Empty(t, user.GetEvents())

	createdEv := events[0]
	assert.Equal(t, fixedEventID, createdEv.ID())
	assert.Equal(t, domain.TopicUserCreated, createdEv.Topic())
	assert.Equal(t, kernel.ID(fixedUserID), createdEv.AggregateID())
	assert.Equal(t, []byte(""), createdEv.Payload())
}

func TestUser_LoadUser_NoEvents(t *testing.T) {
	signed := time.Now().Add(-time.Hour).UTC()
	user := domain.LoadUser(fixedUserID, fixedEmail, fixedUserNum, fixedEncrypted, signed)
	assert.NotNil(t, user)
	assert.Equal(t, fixedUserID, user.ID())
	assert.Equal(t, fixedUserNum, user.Number())
	assert.Equal(t, fixedEmail, user.Email())
	assert.Equal(t, fixedEncrypted, user.PasswordEncrypted())
	assert.Equal(t, signed, user.SignedUpAt())
	assert.Empty(t, user.GetEvents())
}

func TestUser_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	eventIDGenerator.EXPECT().Generate().Return(fixedEventID)

	userIDGenerator := domainMocks.NewMockUserIDGenerator(ctrl)
	userIDGenerator.EXPECT().Generate().Return(fixedUserID)

	userNumberGenerator := domainMocks.NewMockUserNumberGenerator(ctrl)
	userNumberGenerator.EXPECT().Generate().Return(fixedUserNum)

	user, err := domain.CreateUser(
		fixedEmail,
		fixedEncrypted,
		userIDGenerator,
		userNumberGenerator,
		eventIDGenerator,
	)
	require.NoError(t, err)
	_ = user.GetEvents()

	const plainPassword = "plain-pass"

	comparator := domainMocks.NewMockComparator(ctrl)
	comparator.EXPECT().Compare(fixedEncrypted.String(), plainPassword).Return(nil)

	refreshGen := domainMocks.NewMockRefreshTokenGenerator(ctrl)
	refreshGen.EXPECT().Generate().Return(domain.RefreshToken("rt-abc"), nil)

	rt, err := user.Login(plainPassword, comparator, refreshGen)
	require.NoError(t, err)
	require.NotNil(t, rt)
	assert.Equal(t, fixedUserID, rt.UserID())
	assert.Equal(t, domain.RefreshToken("rt-abc"), rt.Token())
	assert.Empty(t, user.GetEvents())
}

func TestUser_Login_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	eventIDGenerator.EXPECT().Generate().Return(fixedEventID)

	userIDGenerator := domainMocks.NewMockUserIDGenerator(ctrl)
	userIDGenerator.EXPECT().Generate().Return(fixedUserID)

	userNumberGenerator := domainMocks.NewMockUserNumberGenerator(ctrl)
	userNumberGenerator.EXPECT().Generate().Return(fixedUserNum)

	user, err := domain.CreateUser(
		fixedEmail,
		fixedEncrypted,
		userIDGenerator,
		userNumberGenerator,
		eventIDGenerator,
	)
	require.NoError(t, err)
	_ = user.GetEvents()

	const plainPassword = "wrong-pass"

	comparator := domainMocks.NewMockComparator(ctrl)
	comparator.EXPECT().Compare(fixedEncrypted.String(), plainPassword).Return(errors.New("password mismatch"))

	refreshGen := domainMocks.NewMockRefreshTokenGenerator(ctrl)
	refreshGen.EXPECT().Generate().Times(0)

	rt, err := user.Login(plainPassword, comparator, refreshGen)
	require.Error(t, err)
	assert.Nil(t, rt)
	assert.Contains(t, err.Error(), "password mismatch")
	assert.Empty(t, user.GetEvents())
}

func TestUser_Login_RefreshTokenGeneratorError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	eventIDGenerator.EXPECT().Generate().Return(fixedEventID)

	userIDGenerator := domainMocks.NewMockUserIDGenerator(ctrl)
	userIDGenerator.EXPECT().Generate().Return(fixedUserID)

	userNumberGenerator := domainMocks.NewMockUserNumberGenerator(ctrl)
	userNumberGenerator.EXPECT().Generate().Return(fixedUserNum)

	user, err := domain.CreateUser(
		fixedEmail,
		fixedEncrypted,
		userIDGenerator,
		userNumberGenerator,
		eventIDGenerator,
	)
	require.NoError(t, err)
	_ = user.GetEvents()

	const plainPassword = "plain-pass"

	comparator := domainMocks.NewMockComparator(ctrl)
	comparator.EXPECT().Compare(fixedEncrypted.String(), plainPassword).Return(nil)

	refreshGen := domainMocks.NewMockRefreshTokenGenerator(ctrl)
	refreshGen.EXPECT().Generate().Return(domain.RefreshToken(""), errors.New("generator failure"))

	rt, err := user.Login(plainPassword, comparator, refreshGen)
	require.Error(t, err)
	assert.Nil(t, rt)
	assert.Contains(t, err.Error(), "generator failure")
	assert.Empty(t, user.GetEvents())
}

func TestUser_Getters(t *testing.T) {
	signed := time.Now().UTC()
	user := domain.LoadUser(fixedUserID, fixedEmail, fixedUserNum, fixedEncrypted, signed)
	assert.Equal(t, fixedUserID, user.ID())
	assert.Equal(t, fixedUserNum, user.Number())
	assert.Equal(t, fixedEmail, user.Email())
	assert.Equal(t, fixedEncrypted, user.PasswordEncrypted())
	assert.Equal(t, signed, user.SignedUpAt())
	assert.Empty(t, user.GetEvents())
}
