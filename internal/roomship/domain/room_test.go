package domain_test

import (
	"gochat/internal/roomship/domain"
	"gochat/internal/roomship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMock "gochat/internal/shared/event/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRoomNumber_Validate(t *testing.T) {
	validNumber := domain.RoomNumber("123456789")
	err := validNumber.Validate()
	require.NoError(t, err)

	invalidNumber := domain.RoomNumber("invalid_number!")
	err = invalidNumber.Validate()
	require.ErrorIs(t, err, myErrors.ErrInvalidNumber)
}

func TestLoadRoom(t *testing.T) {
	start := time.Now().UTC()
	room := domain.LoadRoom(
		fixedRoomID,
		fixedUserID,
		fixedRoomNumber,
		fixedPasswordEncrypted,
		10,
		time.Now().UTC(),
	)
	require.NotNil(t, room)
	assert.WithinDuration(t, start, room.CreatedAt(), timeTolerance)
	assert.Equal(t, fixedRoomID, room.ID())
	assert.Equal(t, fixedUserID, room.OwnerID())
	assert.Equal(t, fixedRoomNumber, room.Number())
	assert.Equal(t, fixedPasswordEncrypted, room.PasswordEncrypted())
	assert.Equal(t, 10, room.MaxMemberCount())
}

func TestCreateRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomIDGenerator := mocks.NewMockRoomIDGenerator(ctrl)
	mockRoomNumberGenerator := mocks.NewMockRoomNumberGenerator(ctrl)
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)

	//正常情况
	mockRoomIDGenerator.EXPECT().Generate().Return(fixedRoomID).Times(1)
	mockRoomNumberGenerator.EXPECT().Generate().Return(fixedRoomNumber).Times(1)
	mockIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1)

	start := time.Now().UTC()
	room, err := domain.CreateRoom(
		fixedUserID,
		20,
		fixedPasswordEncrypted,
		mockRoomIDGenerator,
		mockRoomNumberGenerator,
		mockIDGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, room)
	assert.WithinDuration(t, start, room.CreatedAt(), timeTolerance)
	assert.Equal(t, fixedRoomID, room.ID())
	assert.Equal(t, fixedUserID, room.OwnerID())
	assert.Equal(t, fixedRoomNumber, room.Number())
	assert.Equal(t, fixedPasswordEncrypted, room.PasswordEncrypted())
	assert.Equal(t, 20, room.MaxMemberCount())

	//事件检查
	evs := room.GetEvents()
	require.Len(t, evs, 1)

	createdEv := evs[0]
	assert.Equal(t, domain.TopicRoomCreated, createdEv.Topic())
	assert.Equal(t, fixedRoomID.String(), createdEv.AggregateID().String())
	assert.Equal(t, fixedEventID, createdEv.ID())

	//最大成员数过小
	roomWithInvalidMaxMember, err := domain.CreateRoom(
		fixedUserID,
		1,
		fixedPasswordEncrypted,
		mockRoomIDGenerator,
		mockRoomNumberGenerator,
		mockIDGenerator,
	)
	require.ErrorIs(t, err, domain.ErrInvalidMaxMemberCount)
	require.Nil(t, roomWithInvalidMaxMember)
}
