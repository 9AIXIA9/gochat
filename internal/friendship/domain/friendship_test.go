package domain_test

import (
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/domain/mocks"
	eventMock "gochat/internal/shared/event/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoadFriendship(t *testing.T) {
	start := time.Now().UTC().
		UTC()
	friendship := domain.LoadFriendship(
		fixedFriendshipID,
		fixedUserID,
		fixedToUserID,
		start,
	)
	require.NotNil(t, friendship)
	assert.Equal(t, fixedFriendshipID, friendship.ID())
	assert.Equal(t, fixedUserID, friendship.UserID1())
	assert.Equal(t, fixedToUserID, friendship.UserID2())
	assert.Equal(t, start, friendship.CreatedAt())
}

func TestCreateFriendship(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1)

	mockFriendshipIDGenerator := mocks.NewMockFriendshipIDGenerator(ctrl)
	mockFriendshipIDGenerator.EXPECT().Generate().Return(fixedFriendshipID).Times(1)

	//正常情况
	start := time.Now().UTC().
		UTC()
	friendship, err := domain.CreateFriendship(
		fixedUserID,
		fixedToUserID,
		mockFriendshipIDGenerator,
		mockIDGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, friendship)
	assert.Equal(t, fixedFriendshipID, friendship.ID())
	assert.Equal(t, fixedUserID, friendship.UserID1())
	assert.Equal(t, fixedToUserID, friendship.UserID2())
	assert.WithinDuration(t, start, friendship.CreatedAt(), timeTolerance)

	evs := friendship.GetEvents()
	require.Len(t, evs, 1)
	createdEvent := evs[0]
	assert.Equal(t, domain.TopicFriendshipCreated, createdEvent.Topic())

	//自己添加自己为好友错误
	req, err := domain.CreateFriendship(
		fixedUserID,
		fixedUserID,
		mockFriendshipIDGenerator,
		mockIDGenerator,
	)
	require.ErrorIs(t, err, domain.ErrAddYourselfAsFriend)
	require.Nil(t, req)
}
