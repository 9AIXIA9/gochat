package domain_test

import (
	"gochat/internal/roomship/domain"
	"gochat/internal/roomship/domain/mocks"
	eventMock "gochat/internal/shared/event/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoadRoomship(t *testing.T) {
	start := time.Now().UTC().
		UTC()
	roomship := domain.LoadRoomship(
		fixedRoomshipID,
		fixedRoomID,
		fixedUserID,
		domain.OwnerRole,
		time.Now().UTC(),
	)
	require.NotNil(t, roomship)
	assert.WithinDuration(t, start, roomship.CreatedAt(), timeTolerance)
	assert.Equal(t, fixedRoomshipID, roomship.ID())
	assert.Equal(t, fixedRoomID, roomship.RoomID())
	assert.Equal(t, fixedUserID, roomship.UserID())
	assert.Equal(t, domain.OwnerRole, roomship.Role())
}

func TestCreateRoomship(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomshipIDGenerator := mocks.NewMockRoomshipIDGenerator(ctrl)
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)

	//正常情况
	mockRoomshipIDGenerator.EXPECT().Generate().Return(fixedRoomshipID).Times(1)
	mockIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1)

	start := time.Now().UTC().
		UTC()
	roomship, err := domain.CreateRoomship(
		fixedUserID,
		fixedRoomID,
		domain.MemberRole,
		mockRoomshipIDGenerator,
		mockIDGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, roomship)
	assert.WithinDuration(t, start, roomship.CreatedAt(), timeTolerance)
	assert.Equal(t, fixedUserID, roomship.UserID())
	assert.Equal(t, fixedRoomID, roomship.RoomID())
	assert.Equal(t, domain.MemberRole, roomship.Role())
	assert.Equal(t, fixedRoomshipID, roomship.ID())

	evs := roomship.GetEvents()
	require.Len(t, evs, 1)

	createdEv := evs[0]
	assert.Equal(t, domain.TopicRoomshipCreated, createdEv.Topic())
	assert.Equal(t, fixedRoomshipID.String(), createdEv.AggregateID().String())
	assert.Equal(t, fixedEventID, createdEv.ID())
}
