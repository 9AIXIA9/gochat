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
	fixedRoomIDLeft  kernel.RoomID = "roomship-room-left"
	fixedUserIDLeft  kernel.UserID = "roomship-user-left"
	fixedEventIDLeft event.ID      = "roomship-event-room-left"
)

func TestNewRoomLeftEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedEventIDLeft)

	ev, err := domain.NewRoomLeftEvent(fixedUserIDLeft, fixedRoomIDLeft, idGen)
	require.NoError(t, err)
	require.NotNil(t, ev)
	require.Equal(t, fixedEventIDLeft, ev.ID())
	require.Equal(t, domain.TopicRoomLeft, ev.Topic())
	require.Equal(t, kernel.ID(fixedRoomIDLeft), ev.AggregateID())
	require.Equal(t, fixedUserIDLeft, ev.UserID())
}

func TestToRoomLeftEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedEventIDLeft).AnyTimes()

	payload, _ := json.Marshal(struct{ UserID kernel.UserID }{UserID: fixedUserIDLeft})
	std := event.NewStandardEvent(kernel.ID(fixedRoomIDLeft), domain.TopicRoomLeft, payload, idGen)
	converted, err := domain.ToRoomLeftEvent(std)
	require.NoError(t, err)
	require.NotNil(t, converted)
	require.Equal(t, fixedUserIDLeft, converted.UserID())

	bad := event.NewStandardEvent(kernel.ID(fixedRoomIDLeft), "wrong.topic", payload, idGen)
	converted2, err := domain.ToRoomLeftEvent(bad)
	require.ErrorIs(t, err, myErrors.ErrWrongEventType)
	require.Nil(t, converted2)
}

func TestRoomLeftEvent_MarshalUnmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedEventIDLeft).Times(2)
	ev, err := domain.NewRoomLeftEvent(fixedUserIDLeft, fixedRoomIDLeft, idGen)
	require.NoError(t, err)
	payload, err := ev.Marshal()
	require.NoError(t, err)
	var decoded struct{ UserID kernel.UserID }
	require.NoError(t, json.Unmarshal(payload, &decoded))
	require.Equal(t, fixedUserIDLeft, decoded.UserID)
	std := event.NewStandardEvent(kernel.ID(fixedRoomIDLeft), domain.TopicRoomLeft, payload, idGen)
	converted, err := domain.ToRoomLeftEvent(std)
	require.NoError(t, err)
	require.Equal(t, ev.UserID(), converted.UserID())
}
