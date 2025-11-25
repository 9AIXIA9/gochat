package domain_test

import (
	"errors"
	"gochat/internal/authorization/domain"
	"gochat/internal/authorization/domain/mocks"
	"gochat/internal/shared/event"
	"testing"
	"time"

	eventMock "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const userTimeTolerance = 150 * time.Millisecond

const (
	fixedUserID    kernel.UserID            = "user-123"
	fixedUserNum   kernel.UserNumber        = "10001"
	fixedEventID   event.ID                 = "event-999"
	fixedEmail     kernel.Email             = "test@example.com"
	fixedEncrypted domain.PasswordEncrypted = "encrypted-pass"
	plainPassword                           = "plain-pass"
	wrongPassword                           = "wrong-pass"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	eventIDGenerator.EXPECT().Generate().Return(fixedEventID)

	userIDGenerator := mocks.NewMockUserIDGenerator(ctrl)
	userIDGenerator.EXPECT().Generate().Return(fixedUserID)

	userNumberGenerator := mocks.NewMockUserNumberGenerator(ctrl)
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
	// Use creation start time rather than another now() to reduce race window
	assert.WithinDuration(t, start, user.SignedUpAt(), userTimeTolerance)

	events := user.GetEvents()
	assert.Len(t, events, 1)
	assert.Empty(t, user.GetEvents()) // drained

	createdEv := events[0]
	assert.Equal(t, fixedEventID, createdEv.ID())
	assert.Equal(t, domain.TopicUserCreated, createdEv.Topic())
	assert.Equal(t, kernel.ID(fixedUserID), createdEv.AggregateID())
	assert.Equal(t, []byte(""), createdEv.Payload())
}

func TestUser_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Table-driven scenarios
	type testCase struct {
		name              string
		password          domain.Password
		setupComparator   func(*gomock.Controller) domain.Comparator
		setupRefreshGen   func(*gomock.Controller) domain.RefreshTokenGenerator
		expectErrContains string
		assertToken       func(t *testing.T, rt *domain.RefreshTokenEntity)
	}

	cases := []testCase{
		{
			name:     "success",
			password: domain.Password(plainPassword),
			setupComparator: func(c *gomock.Controller) domain.Comparator {
				m := mocks.NewMockComparator(c)
				m.EXPECT().Compare(fixedEncrypted.String(), plainPassword).Return(nil)
				return m
			},
			setupRefreshGen: func(c *gomock.Controller) domain.RefreshTokenGenerator {
				m := mocks.NewMockRefreshTokenGenerator(c)
				m.EXPECT().Generate().Return(domain.RefreshToken("rt-abc"), nil)
				return m
			},
			assertToken: func(t *testing.T, rt *domain.RefreshTokenEntity) {
				require.NotNil(t, rt)
				assert.Equal(t, fixedUserID, rt.UserID())
				assert.Equal(t, domain.RefreshToken("rt-abc"), rt.Token())
			},
		},
		{
			name:     "wrong password",
			password: domain.Password(wrongPassword),
			setupComparator: func(c *gomock.Controller) domain.Comparator {
				m := mocks.NewMockComparator(c)
				m.EXPECT().Compare(fixedEncrypted.String(), wrongPassword).Return(errors.New("password mismatch"))
				return m
			},
			setupRefreshGen: func(c *gomock.Controller) domain.RefreshTokenGenerator {
				m := mocks.NewMockRefreshTokenGenerator(c)
				// Expect no Generate call
				m.EXPECT().Generate().Times(0)
				return m
			},
			expectErrContains: "password mismatch",
			assertToken: func(t *testing.T, rt *domain.RefreshTokenEntity) {
				assert.Nil(t, rt)
			},
		},
		{
			name:     "refresh token generator error",
			password: domain.Password(plainPassword),
			setupComparator: func(c *gomock.Controller) domain.Comparator {
				m := mocks.NewMockComparator(c)
				m.EXPECT().Compare(fixedEncrypted.String(), plainPassword).Return(nil)
				return m
			},
			setupRefreshGen: func(c *gomock.Controller) domain.RefreshTokenGenerator {
				m := mocks.NewMockRefreshTokenGenerator(c)
				m.EXPECT().Generate().Return(domain.RefreshToken(""), errors.New("generator failure"))
				return m
			},
			expectErrContains: "generator failure",
			assertToken: func(t *testing.T, rt *domain.RefreshTokenEntity) {
				assert.Nil(t, rt)
			},
		},
	}

	for _, tc := range cases {
		// use t.Run for independent subtests
		c := tc // capture
		user := createUser(t, ctrl)
		// drain creation event to isolate login effects
		_ = user.GetEvents()

		comp := c.setupComparator(ctrl)
		refGen := c.setupRefreshGen(ctrl)

		rt, err := user.Login(c.password, comp, refGen)
		if c.expectErrContains != "" {
			require.Error(t, err, c.name)
			assert.Contains(t, err.Error(), c.expectErrContains)
		} else {
			require.NoError(t, err, c.name)
		}
		c.assertToken(t, rt)
		assert.Empty(t, user.GetEvents()) // login should not emit domain events
	}
}

func createUser(t *testing.T, ctrl *gomock.Controller) *domain.User {
	// helper to create user with deterministic generators
	// events will contain one user_created which tests may drain.
	t.Helper()
	eventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	eventIDGenerator.EXPECT().Generate().Return(fixedEventID)

	userIDGenerator := mocks.NewMockUserIDGenerator(ctrl)
	userIDGenerator.EXPECT().Generate().Return(fixedUserID)

	userNumberGenerator := mocks.NewMockUserNumberGenerator(ctrl)
	userNumberGenerator.EXPECT().Generate().Return(fixedUserNum)

	user, err := domain.CreateUser(
		fixedEmail,
		fixedEncrypted,
		userIDGenerator,
		userNumberGenerator,
		eventIDGenerator,
	)
	require.NoError(t, err)
	return user
}
