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
	fixedFriendshipUserEventID     event.ID          = "friendship-user-ev-1"
	fixedFriendshipUserEventUserID kernel.UserID     = "friendship-user-777"
	fixedFriendshipUserEventNumber kernel.UserNumber = "70001"
)

func TestFriendshipUserCreatedEvent_NewUserCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedFriendshipUserEventID)

	ev, err := domain.NewUserCreatedEvent(fixedFriendshipUserEventUserID, fixedFriendshipUserEventNumber, idGen)
	require.NoError(t, err)
	require.NotNil(t, ev)
	require.Equal(t, fixedFriendshipUserEventID, ev.ID())
	require.Equal(t, domain.TopicUserCreated, ev.Topic())
	require.Equal(t, kernel.ID(fixedFriendshipUserEventUserID), ev.AggregateID())
	require.Equal(t, fixedFriendshipUserEventNumber, ev.Number())
	// payload should be json with number
	require.NotEmpty(t, ev.Payload())
}

func TestFriendshipUserCreatedEvent_ToUserCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedFriendshipUserEventID).AnyTimes()

	std := event.NewStandardEvent(kernel.ID(fixedFriendshipUserEventUserID), domain.TopicUserCreated, []byte("{\"Number\":\""+fixedFriendshipUserEventNumber.String()+"\"}"), idGen)
	parsed, err := domain.ToUserCreatedEvent(std)
	require.NoError(t, err)
	require.NotNil(t, parsed)
	require.Equal(t, fixedFriendshipUserEventNumber, parsed.Number())

	wrong := event.NewStandardEvent(kernel.ID(fixedFriendshipUserEventUserID), "wrong.topic", []byte("{}"), idGen)
	parsed2, err := domain.ToUserCreatedEvent(wrong)
	require.ErrorIs(t, err, errors.ErrWrongEventType)
	require.Nil(t, parsed2)
}

func TestFriendshipUserCreatedEvent_MarshalUnmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedFriendshipUserEventID).Times(2)
	ev, err := domain.NewUserCreatedEvent(fixedFriendshipUserEventUserID, fixedFriendshipUserEventNumber, idGen)
	require.NoError(t, err)
	payload, err := ev.Marshal()
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	std := event.NewStandardEvent(kernel.ID(fixedFriendshipUserEventUserID), domain.TopicUserCreated, payload, idGen)
	converted, err := domain.ToUserCreatedEvent(std)
	require.NoError(t, err)
	require.Equal(t, ev.Number(), converted.Number())
}
