package domain_test

import (
	"testing"
	"time"

	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	eventMocks "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"
	"gochat/internal/social/domain/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedOwnerID    kernel.UserID     = "owner-1"
	fixedUserID2    kernel.UserID     = "user-2"
	fixedUserID3    kernel.UserID     = "user-3"
	fixedRoomID     kernel.RoomID     = "room-999"
	fixedRoomNumber kernel.RoomNumber = "123456"
	fixedEventID    event.ID          = "event-abc"
)

func TestRoomOption_Validate(t *testing.T) {
	opt := &domain.RoomOption{MaxMemberCount: 5}
	require.NoError(t, opt.Validate())
	bad := &domain.RoomOption{MaxMemberCount: 1}
	require.ErrorIs(t, bad.Validate(), myErrors.ErrLessThanMin)
}

func TestCreateRoom_DefaultOption(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomIDGen := mocks.NewMockRoomIDGenerator(ctrl)
	roomIDGen.EXPECT().Generate().Return(fixedRoomID)
	roomNumberGen := mocks.NewMockRoomNumberGenerator(ctrl)
	roomNumberGen.EXPECT().Generate().Return(fixedRoomNumber)
	eventIDGen := eventMocks.NewMockIDGenerator(ctrl)
	eventIDGen.EXPECT().Generate().Return(fixedEventID) // room created event

	r, err := domain.CreateRoom(fixedOwnerID, roomIDGen, roomNumberGen, eventIDGen)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, fixedRoomID, r.ID())
	assert.Equal(t, fixedRoomNumber, r.Number())
	assert.Equal(t, fixedOwnerID, r.OwnerID())
	assert.Equal(t, 1, r.MemberCount())
	assert.Equal(t, 20, r.MaxMemberCount()) // default
	assert.WithinDuration(t, time.Now().UTC(), r.CreatedAt(), 200*time.Millisecond)

	events := r.GetEvents()
	require.Len(t, events, 1)
	assert.Equal(t, domain.TopicRoomCreated, events[0].Topic())
	assert.Equal(t, kernel.ID(fixedRoomID), events[0].AggregateID())
}

func TestCreateRoom_WithOptionInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	opt := &domain.RoomOption{MaxMemberCount: 1}
	roomIDGen := mocks.NewMockRoomIDGenerator(ctrl)
	roomNumberGen := mocks.NewMockRoomNumberGenerator(ctrl)
	eventIDGen := eventMocks.NewMockIDGenerator(ctrl)
	// no Generate calls expected due to validation failure

	r, err := domain.CreateRoom(fixedOwnerID, roomIDGen, roomNumberGen, eventIDGen, opt)
	require.ErrorIs(t, err, myErrors.ErrLessThanMin)
	require.Nil(t, r)
}

func TestRoom_AddMember_Success_NoPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomIDGen := mocks.NewMockRoomIDGenerator(ctrl)
	roomIDGen.EXPECT().Generate().Return(fixedRoomID)
	roomNumberGen := mocks.NewMockRoomNumberGenerator(ctrl)
	roomNumberGen.EXPECT().Generate().Return(fixedRoomNumber)
	eventIDGen := eventMocks.NewMockIDGenerator(ctrl)
	// first for create, second for join
	gomock.InOrder(
		eventIDGen.EXPECT().Generate().Return(fixedEventID),
		eventIDGen.EXPECT().Generate().Return(fixedEventID),
	)

	r, err := domain.CreateRoom(fixedOwnerID, roomIDGen, roomNumberGen, eventIDGen)
	require.NoError(t, err)

	err = r.AddMember(fixedUserID2, "", nil, eventIDGen)
	require.NoError(t, err)
	assert.Equal(t, 2, r.MemberCount())

	events := r.GetEvents()
	require.Len(t, events, 2)
	assert.Equal(t, domain.TopicRoomJoined, events[1].Topic())
}

func TestRoom_AddMember_WithPassword_SuccessAndComparatorError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	encryptedPass := domain.PasswordEncrypted("enc-pass")
	roomIDGen := mocks.NewMockRoomIDGenerator(ctrl)
	roomIDGen.EXPECT().Generate().Return(fixedRoomID)
	roomNumberGen := mocks.NewMockRoomNumberGenerator(ctrl)
	roomNumberGen.EXPECT().Generate().Return(fixedRoomNumber)

	// success path event id sequence: create + join
	eventIDGenSuccess := eventMocks.NewMockIDGenerator(ctrl)
	gomock.InOrder(
		eventIDGenSuccess.EXPECT().Generate().Return(fixedEventID),
		eventIDGenSuccess.EXPECT().Generate().Return(fixedEventID),
	)

	opt := &domain.RoomOption{MaxMemberCount: 5, PasswordEncrypted: encryptedPass}
	r, err := domain.CreateRoom(fixedOwnerID, roomIDGen, roomNumberGen, eventIDGenSuccess, opt)
	require.NoError(t, err)

	comp := mocks.NewMockComparator(ctrl)
	comp.EXPECT().Compare(encryptedPass.String(), "raw-pass").Return(nil)
	err = r.AddMember(fixedUserID2, "raw-pass", comp, eventIDGenSuccess)
	require.NoError(t, err)
	assert.Equal(t, 2, r.MemberCount())

	// comparator error path (no extra event expected)
	roomIDGenFail := mocks.NewMockRoomIDGenerator(ctrl)
	roomIDGenFail.EXPECT().Generate().Return(fixedRoomID)

	roomNumberGenFail := mocks.NewMockRoomNumberGenerator(ctrl)
	roomNumberGenFail.EXPECT().Generate().Return(fixedRoomNumber)

	eventIDGenFail := eventMocks.NewMockIDGenerator(ctrl)
	// only create room event
	gomock.InOrder(
		eventIDGenFail.EXPECT().Generate().Return(fixedEventID),
	)
	r2, err := domain.CreateRoom(fixedOwnerID, roomIDGenFail, roomNumberGenFail, eventIDGenFail, opt)
	require.NoError(t, err)

	compErr := mocks.NewMockComparator(ctrl)
	compErr.EXPECT().Compare(encryptedPass.String(), "bad-pass").Return(myErrors.ErrInvalidCredential)
	err = r2.AddMember(fixedUserID2, "bad-pass", compErr, eventIDGenFail)
	require.ErrorIs(t, err, myErrors.ErrInvalidCredential)
	assert.Equal(t, 1, r2.MemberCount())
	assert.Len(t, r2.GetEvents(), 1) // only room created
}

func TestRoom_AddMember_ExceedMaxAndDuplicate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomIDGen := mocks.NewMockRoomIDGenerator(ctrl)
	roomIDGen.EXPECT().Generate().Return(fixedRoomID)
	roomNumberGen := mocks.NewMockRoomNumberGenerator(ctrl)
	roomNumberGen.EXPECT().Generate().Return(fixedRoomNumber)
	eventIDGen := eventMocks.NewMockIDGenerator(ctrl)
	// create + join second member
	gomock.InOrder(
		eventIDGen.EXPECT().Generate().Return(fixedEventID),
		eventIDGen.EXPECT().Generate().Return(fixedEventID),
	)

	opt := &domain.RoomOption{MaxMemberCount: 2}
	r, err := domain.CreateRoom(fixedOwnerID, roomIDGen, roomNumberGen, eventIDGen, opt)
	require.NoError(t, err)

	err = r.AddMember(fixedUserID2, "", nil, eventIDGen)
	require.NoError(t, err)
	assert.Equal(t, 2, r.MemberCount())

	// exceed max
	err = r.AddMember(fixedUserID3, "", nil, eventIDGen)
	require.ErrorIs(t, err, myErrors.ErrExceedMaxValue)
	assert.Equal(t, 2, r.MemberCount())

	// duplicate add (owner again) should be no-op and no new event
	eventIDGenDup := eventMocks.NewMockIDGenerator(ctrl)
	// only create room event expected
	roomIDGen2 := mocks.NewMockRoomIDGenerator(ctrl)
	roomIDGen2.EXPECT().Generate().Return(fixedRoomID)
	roomNumberGen2 := mocks.NewMockRoomNumberGenerator(ctrl)
	roomNumberGen2.EXPECT().Generate().Return(fixedRoomNumber)
	gomock.InOrder(
		eventIDGenDup.EXPECT().Generate().Return(fixedEventID),
	)
	r2, err := domain.CreateRoom(fixedOwnerID, roomIDGen2, roomNumberGen2, eventIDGenDup, opt)
	require.NoError(t, err)
	err = r2.AddMember(fixedOwnerID, "", nil, eventIDGenDup)
	require.NoError(t, err)
	assert.Equal(t, 1, r2.MemberCount())
	assert.Len(t, r2.GetEvents(), 1)
}

func TestRoom_DeleteMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomIDGen := mocks.NewMockRoomIDGenerator(ctrl)
	roomIDGen.EXPECT().Generate().Return(fixedRoomID)
	roomNumberGen := mocks.NewMockRoomNumberGenerator(ctrl)
	roomNumberGen.EXPECT().Generate().Return(fixedRoomNumber)

	// create + join + leave event id expectation sequence
	eventIDGen := eventMocks.NewMockIDGenerator(ctrl)
	gomock.InOrder(
		eventIDGen.EXPECT().Generate().Return(fixedEventID), // create
		eventIDGen.EXPECT().Generate().Return(fixedEventID), // join
		eventIDGen.EXPECT().Generate().Return(fixedEventID), // leave
	)

	r, err := domain.CreateRoom(fixedOwnerID, roomIDGen, roomNumberGen, eventIDGen)
	require.NoError(t, err)
	err = r.AddMember(fixedUserID2, "", nil, eventIDGen)
	require.NoError(t, err)

	// delete success
	err = r.DeleteMember(fixedUserID2, eventIDGen)
	require.NoError(t, err)
	assert.Equal(t, 1, r.MemberCount())

	events := r.GetEvents()
	require.Len(t, events, 3)
	assert.Equal(t, domain.TopicRoomLeft, events[2].Topic())
	// owner cant leave
	eventIDGenOwner := eventMocks.NewMockIDGenerator(ctrl)
	gomock.InOrder(
		eventIDGenOwner.EXPECT().Generate().Return(fixedEventID), // only create, no leave event
	)
	roomIDGen2 := mocks.NewMockRoomIDGenerator(ctrl)
	roomIDGen2.EXPECT().Generate().Return(fixedRoomID)
	roomNumberGen2 := mocks.NewMockRoomNumberGenerator(ctrl)
	roomNumberGen2.EXPECT().Generate().Return(fixedRoomNumber)
	r2, err := domain.CreateRoom(fixedOwnerID, roomIDGen2, roomNumberGen2, eventIDGenOwner)
	require.NoError(t, err)
	err = r2.DeleteMember(fixedOwnerID, eventIDGenOwner)
	require.ErrorIs(t, err, myErrors.ErrOwnerCantLeave)
	assert.Len(t, r2.GetEvents(), 1)

	// non-member delete (no error, no event)
	eventIDGenNon := eventMocks.NewMockIDGenerator(ctrl)
	gomock.InOrder(
		eventIDGenNon.EXPECT().Generate().Return(fixedEventID),
	)
	roomIDGen3 := mocks.NewMockRoomIDGenerator(ctrl)
	roomIDGen3.EXPECT().Generate().Return(fixedRoomID)
	roomNumberGen3 := mocks.NewMockRoomNumberGenerator(ctrl)
	roomNumberGen3.EXPECT().Generate().Return(fixedRoomNumber)
	r3, err := domain.CreateRoom(fixedOwnerID, roomIDGen3, roomNumberGen3, eventIDGenNon)
	require.NoError(t, err)
	err = r3.DeleteMember(fixedUserID2, eventIDGenNon)
	require.NoError(t, err)
	assert.Len(t, r3.GetEvents(), 1)
}

func TestRoom_LoadRoom(t *testing.T) {
	createdAt := time.Now().Add(-1 * time.Hour).UTC()
	members := []kernel.UserID{fixedOwnerID, fixedUserID2}

	room := domain.LoadRoom(
		fixedRoomID,
		fixedOwnerID,
		fixedRoomNumber,
		fixedEncryptedPassword,
		members,
		10,
		createdAt,
	)
	assert.Equal(t, fixedRoomID, room.ID())
	assert.Equal(t, fixedRoomNumber, room.Number())
	assert.Equal(t, fixedOwnerID, room.OwnerID())
	assert.Equal(t, 2, room.MemberCount())
	assert.Equal(t, 10, room.MaxMemberCount())
	assert.Equal(t, createdAt, room.CreatedAt())
	assert.Equal(t, fixedEncryptedPassword, room.PasswordEncrypted().String())
	assert.ElementsMatch(t, members, room.Members())
}
