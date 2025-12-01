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
	fixedUserID                             kernel.UserID = "user-123"
	fixedFriendID                           kernel.UserID = "friend-user-456"
	fixedFriendshipFriendshipCreatedEventID event.ID      = "ev-friendship-created-1"
)

func TestNewFriendshipCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedFriendshipFriendshipCreatedEventID)

	ev, err := domain.NewFriendshipCreatedEvent(
		fixedUserID,
		fixedFriendID,
		idGen,
	)
	require.NoError(t, err)
	require.NotNil(t, ev)
	require.Equal(t, fixedFriendshipFriendshipCreatedEventID, ev.ID())
	require.Equal(t, domain.TopicFriendshipCreated, ev.Topic())
	require.Equal(t, kernel.ID(fixedUserID), ev.AggregateID())
	require.Equal(t, fixedFriendID, ev.FriendID())
	require.NotEmpty(t, ev.Payload())
}

func TestToFriendshipCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedFriendshipFriendshipCreatedEventID).AnyTimes()

	std := event.NewStandardEvent(
		kernel.ID(fixedFriendshipUserEventUserID),
		domain.TopicFriendshipCreated,
		[]byte("{\"FriendID\":\""+fixedFriendID.String()+"\"}"),
		idGen,
	)
	parsed, err := domain.ToFriendshipCreatedEvent(std)
	require.NoError(t, err)
	require.NotNil(t, parsed)
	require.Equal(t, fixedFriendID, parsed.FriendID())

	wrong := event.NewStandardEvent(kernel.ID(fixedFriendshipUserEventUserID), "wrong.topic", []byte("{}"), idGen)
	parsed2, err := domain.ToFriendshipCreatedEvent(wrong)
	require.ErrorIs(t, err, errors.ErrWrongEventType)
	require.Nil(t, parsed2)
}

func TestFriendshipFriendshipCreatedEvent_MarshalUnmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedFriendshipFriendshipCreatedEventID).Times(2)
	ev, err := domain.NewFriendshipCreatedEvent(fixedFriendshipUserEventUserID, fixedFriendID, idGen)
	require.NoError(t, err)
	payload, err := ev.Marshal()
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	std := event.NewStandardEvent(kernel.ID(fixedFriendshipUserEventUserID), domain.TopicFriendshipCreated, payload, idGen)
	converted, err := domain.ToFriendshipCreatedEvent(std)
	require.NoError(t, err)
	require.Equal(t, ev.FriendID(), converted.FriendID())
}
