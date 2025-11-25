package domain_test

import (
	"encoding/json"
	"testing"

	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	eventMocks "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedRoomIDJoined  kernel.RoomID = "social-room-joined"
	fixedUserIDJoined  kernel.UserID = "social-user-joined"
	fixedEventIDJoined event.ID      = "social-event-room-joined"
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
