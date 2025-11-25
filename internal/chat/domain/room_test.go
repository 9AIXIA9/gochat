package domain_test

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fixedRoomID     kernel.RoomID     = "room-abc"
	fixedRoomNumber kernel.RoomNumber = "30001"
	fixedOwnerID    kernel.UserID     = "owner-999"
	fixedMember1ID  kernel.UserID     = "member-1"
	fixedMember2ID  kernel.UserID     = "member-2"
)

func TestRoom_CreateRoom(t *testing.T) {
	r := domain.CreateRoom(fixedRoomID, fixedRoomNumber, fixedOwnerID)
	require.NotNil(t, r)
	assert.Equal(t, fixedRoomID, r.ID())
	assert.Equal(t, fixedRoomNumber, r.Number())
	assert.Equal(t, []kernel.UserID{fixedOwnerID}, r.Members())
	assert.True(t, r.IsMember(fixedOwnerID))
}

func TestRoom_LoadRoom(t *testing.T) {
	members := []kernel.UserID{fixedOwnerID, fixedMember1ID}
	r := domain.LoadRoom(fixedRoomID, fixedRoomNumber, members)
	require.NotNil(t, r)
	assert.Equal(t, fixedRoomID, r.ID())
	assert.Equal(t, fixedRoomNumber, r.Number())
	assert.ElementsMatch(t, members, r.Members())
}

func TestRoom_AddMember(t *testing.T) {
	r := domain.CreateRoom(fixedRoomID, fixedRoomNumber, fixedOwnerID)
	// add new member
	r.AddMember(fixedMember1ID)
	assert.ElementsMatch(t, []kernel.UserID{fixedOwnerID, fixedMember1ID}, r.Members())
	// add duplicate should not change
	r.AddMember(fixedMember1ID)
	assert.ElementsMatch(t, []kernel.UserID{fixedOwnerID, fixedMember1ID}, r.Members())
	// add second member
	r.AddMember(fixedMember2ID)
	assert.ElementsMatch(t, []kernel.UserID{fixedOwnerID, fixedMember1ID, fixedMember2ID}, r.Members())
}

func TestRoom_DeleteMember(t *testing.T) {
	r := domain.LoadRoom(fixedRoomID, fixedRoomNumber, []kernel.UserID{fixedOwnerID, fixedMember1ID, fixedMember2ID})
	// delete existing
	r.DeleteMember(fixedMember1ID)
	assert.ElementsMatch(t, []kernel.UserID{fixedOwnerID, fixedMember2ID}, r.Members())
	// delete non-existent should no-op
	r.DeleteMember(fixedMember1ID)
	assert.ElementsMatch(t, []kernel.UserID{fixedOwnerID, fixedMember2ID}, r.Members())
	// delete last
	r.DeleteMember(fixedMember2ID)
	assert.ElementsMatch(t, []kernel.UserID{fixedOwnerID}, r.Members())
}
