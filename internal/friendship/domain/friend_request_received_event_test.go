package domain_test

import (
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	eventMocks "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedRequestID                       kernel.OperationID = "friend-request-456"
	fixedFriendshipFriendReceivedEventID event.ID           = "ev-friendship-created-1"
)

func TestNewFriendRequestReceivedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedFriendshipFriendReceivedEventID)

	ev, err := domain.NewFriendRequestReceivedEvent(
		fixedUserID,
		fixedRequestID,
		idGen,
	)
	require.NoError(t, err)
	require.NotNil(t, ev)
	require.Equal(t, fixedFriendshipFriendReceivedEventID, ev.ID())
	require.Equal(t, domain.TopicFriendRequestReceived, ev.Topic())
	require.Equal(t, kernel.ID(fixedUserID), ev.AggregateID())
	require.Equal(t, fixedRequestID, ev.RequestID())
	require.NotEmpty(t, ev.Payload())
}

func TestToFriendRequestReceivedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedFriendshipFriendReceivedEventID).AnyTimes()

	std := event.NewStandardEvent(
		kernel.ID(fixedFriendshipUserEventUserID),
		domain.TopicFriendRequestReceived,
		[]byte("{\"RequestID\":\""+fixedRequestID.String()+"\"}"),
		idGen,
	)
	parsed, err := domain.ToFriendRequestReceivedEvent(std)
	require.NoError(t, err)
	require.NotNil(t, parsed)
	require.Equal(t, fixedRequestID, parsed.RequestID())

	wrong := event.NewStandardEvent(kernel.ID(fixedFriendshipUserEventUserID), "wrong.topic", []byte("{}"), idGen)
	parsed2, err := domain.ToFriendRequestReceivedEvent(wrong)
	require.ErrorIs(t, err, errors.ErrWrongEventType)
	require.Nil(t, parsed2)
}

func TestFriendshipFriendRequestReceivedEvent_MarshalUnmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedFriendshipFriendReceivedEventID).Times(2)
	ev, err := domain.NewFriendRequestReceivedEvent(fixedFriendshipUserEventUserID, fixedRequestID, idGen)
	require.NoError(t, err)
	payload, err := ev.Marshal()
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	std := event.NewStandardEvent(kernel.ID(fixedFriendshipUserEventUserID), domain.TopicFriendRequestReceived, payload, idGen)
	converted, err := domain.ToFriendRequestReceivedEvent(std)
	require.NoError(t, err)
	require.Equal(t, ev.RequestID(), converted.RequestID())
}
