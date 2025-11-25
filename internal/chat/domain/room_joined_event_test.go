package domain_test

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	eventMocks "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedRoomJoinedEventID event.ID      = "room-joined-ev-1"
	fixedRoomJoinedRoomID  kernel.RoomID = "room-join-321"
	fixedRoomJoinedUserID  kernel.UserID = "user-join-abc"
)

func TestRoomJoinedEvent_NewRoomJoinedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedRoomJoinedEventID)

	ev, err := domain.NewRoomJoinedEvent(fixedRoomJoinedRoomID, fixedRoomJoinedUserID, idGen)
	require.NoError(t, err)
	require.NotNil(t, ev)
	require.Equal(t, fixedRoomJoinedEventID, ev.ID())
	require.Equal(t, domain.TopicRoomJoined, ev.Topic())
	require.Equal(t, kernel.ID(fixedRoomJoinedRoomID), ev.AggregateID())
	require.Equal(t, fixedRoomJoinedUserID, ev.UserID())
}

func TestRoomJoinedEvent_ToRoomJoinedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedRoomJoinedEventID).AnyTimes()

	payload := []byte("{\"UserID\":\"" + fixedRoomJoinedUserID.String() + "\"}")
	std := event.NewStandardEvent(kernel.ID(fixedRoomJoinedRoomID), domain.TopicRoomJoined, payload, idGen)
	parsed, err := domain.ToRoomJoinedEvent(std)
	require.NoError(t, err)
	require.NotNil(t, parsed)
	require.Equal(t, fixedRoomJoinedUserID, parsed.UserID())

	wrong := event.NewStandardEvent(kernel.ID(fixedRoomJoinedRoomID), "wrong.topic", payload, idGen)
	parsed2, err := domain.ToRoomJoinedEvent(wrong)
	require.ErrorIs(t, err, errors.ErrWrongEventType)
	require.Nil(t, parsed2)
}
