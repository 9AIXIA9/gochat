package domain_test

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoom_CreateRoom(t *testing.T) {
	r := domain.CreateRoom(fixedRoomID, fixedUserID)
	require.NotNil(t, r)
	assert.Equal(t, fixedRoomID, r.ID())
	assert.Equal(t, []kernel.UserID{fixedUserID}, r.Members())
	assert.True(t, r.IsMember(fixedUserID))
}

func TestRoom_LoadRoom(t *testing.T) {
	members := []kernel.UserID{fixedUserID, fixedMemberID1}
	r := domain.LoadRoom(fixedRoomID, members)
	require.NotNil(t, r)
	assert.Equal(t, fixedRoomID, r.ID())
	assert.ElementsMatch(t, members, r.Members())
}

func TestRoom_AddMember(t *testing.T) {
	r := domain.CreateRoom(fixedRoomID, fixedUserID)
	// add new member
	r.AddMember(fixedMemberID1)
	assert.ElementsMatch(t, []kernel.UserID{fixedUserID, fixedMemberID1}, r.Members())
	// add duplicate should not change
	r.AddMember(fixedMemberID1)
	assert.ElementsMatch(t, []kernel.UserID{fixedUserID, fixedMemberID1}, r.Members())
	// add second member
	r.AddMember(fixedMemberID2)
	assert.ElementsMatch(t, []kernel.UserID{fixedUserID, fixedMemberID1, fixedMemberID2}, r.Members())
}

func TestRoom_DeleteMember(t *testing.T) {
	r := domain.LoadRoom(fixedRoomID, []kernel.UserID{fixedUserID, fixedMemberID1, fixedMemberID2})
	// delete existing
	r.DeleteMember(fixedMemberID1)
	assert.ElementsMatch(t, []kernel.UserID{fixedUserID, fixedMemberID2}, r.Members())
	// delete non-existent should no-op
	r.DeleteMember(fixedMemberID1)
	assert.ElementsMatch(t, []kernel.UserID{fixedUserID, fixedMemberID2}, r.Members())
	// delete last
	r.DeleteMember(fixedMemberID2)
	assert.ElementsMatch(t, []kernel.UserID{fixedUserID}, r.Members())
}
