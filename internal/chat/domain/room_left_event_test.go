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
	fixedRoomLeftEventID event.ID      = "room-left-ev-1"
	fixedRoomLeftRoomID  kernel.RoomID = "room-left-987"
	fixedRoomLeftUserID  kernel.UserID = "user-left-xyz"
)

func TestRoomLeftEvent_NewRoomLeftEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedRoomLeftEventID)

	ev, err := domain.NewRoomLeftEvent(fixedRoomLeftRoomID, fixedRoomLeftUserID, idGen)
	require.NoError(t, err)
	require.NotNil(t, ev)
	require.Equal(t, fixedRoomLeftEventID, ev.ID())
	require.Equal(t, domain.TopicRoomLeft, ev.Topic())
	require.Equal(t, kernel.ID(fixedRoomLeftRoomID), ev.AggregateID())
	require.Equal(t, fixedRoomLeftUserID, ev.UserID())
}

func TestRoomLeftEvent_ToRoomLeftEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedRoomLeftEventID).AnyTimes()

	payload := []byte("{\"UserID\":\"" + fixedRoomLeftUserID.String() + "\"}")
	std := event.NewStandardEvent(kernel.ID(fixedRoomLeftRoomID), domain.TopicRoomLeft, payload, idGen)
	parsed, err := domain.ToRoomLeftEvent(std)
	require.NoError(t, err)
	require.NotNil(t, parsed)
	require.Equal(t, fixedRoomLeftUserID, parsed.UserID())

	wrong := event.NewStandardEvent(kernel.ID(fixedRoomLeftRoomID), "wrong.topic", payload, idGen)
	parsed2, err := domain.ToRoomLeftEvent(wrong)
	require.ErrorIs(t, err, errors.ErrWrongEventType)
	require.Nil(t, parsed2)
}
