package domain_test

import (
	"encoding/json"
	"testing"

	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	eventMocks "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedRoomIDJoined  kernel.RoomID = "roomship-room-joined"
	fixedUserIDJoined  kernel.UserID = "roomship-user-joined"
	fixedEventIDJoined event.ID      = "roomship-event-room-joined"
)

func TestNewRoomJoinedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedEventIDJoined)

	ev, err := domain.NewRoomJoinedEvent(fixedUserIDJoined, fixedRoomIDJoined, idGen)
	require.NoError(t, err)
	require.NotNil(t, ev)
	require.Equal(t, fixedEventIDJoined, ev.ID())
	require.Equal(t, domain.TopicRoomJoined, ev.Topic())
	require.Equal(t, kernel.ID(fixedRoomIDJoined), ev.AggregateID())
	require.Equal(t, fixedUserIDJoined, ev.UserID())
}

func TestToRoomJoinedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedEventIDJoined).AnyTimes()

	// from standard event with payload
	payload, _ := json.Marshal(struct{ UserID kernel.UserID }{UserID: fixedUserIDJoined})
	std := event.NewStandardEvent(kernel.ID(fixedRoomIDJoined), domain.TopicRoomJoined, payload, idGen)
	converted, err := domain.ToRoomJoinedEvent(std)
	require.NoError(t, err)
	require.NotNil(t, converted)
	require.Equal(t, fixedUserIDJoined, converted.UserID())

	// wrong topic
	bad := event.NewStandardEvent(kernel.ID(fixedRoomIDJoined), "wrong.topic", payload, idGen)
	converted2, err := domain.ToRoomJoinedEvent(bad)
	require.ErrorIs(t, err, myErrors.ErrWrongEventType)
	require.Nil(t, converted2)
}

func TestRoomJoinedEvent_MarshalUnmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedEventIDJoined).Times(2)
	ev, err := domain.NewRoomJoinedEvent(fixedUserIDJoined, fixedRoomIDJoined, idGen)
	require.NoError(t, err)
	payload, err := ev.Marshal()
	require.NoError(t, err)
	// decode JSON and verify
	var decoded struct{ UserID kernel.UserID }
	require.NoError(t, json.Unmarshal(payload, &decoded))
	require.Equal(t, fixedUserIDJoined, decoded.UserID)
	// round trip via standard event
	std := event.NewStandardEvent(kernel.ID(fixedRoomIDJoined), domain.TopicRoomJoined, payload, idGen)
	converted, err := domain.ToRoomJoinedEvent(std)
	require.NoError(t, err)
	require.Equal(t, ev.UserID(), converted.UserID())
}
