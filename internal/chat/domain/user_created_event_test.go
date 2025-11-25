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
	fixedChatUserEventID     event.ID          = "chat-user-ev-1"
	fixedChatUserEventUserID kernel.UserID     = "chat-user-777"
	fixedChatUserEventNumber kernel.UserNumber = "70001"
)

func TestChatUserCreatedEvent_NewUserCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedChatUserEventID)

	ev, err := domain.NewUserCreatedEvent(fixedChatUserEventUserID, fixedChatUserEventNumber, idGen)
	require.NoError(t, err)
	require.NotNil(t, ev)
	require.Equal(t, fixedChatUserEventID, ev.ID())
	require.Equal(t, domain.TopicUserCreated, ev.Topic())
	require.Equal(t, kernel.ID(fixedChatUserEventUserID), ev.AggregateID())
	require.Equal(t, fixedChatUserEventNumber, ev.Number())
	// payload should be json with number
	require.NotEmpty(t, ev.Payload())
}

func TestChatUserCreatedEvent_ToUserCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedChatUserEventID).AnyTimes()

	std := event.NewStandardEvent(kernel.ID(fixedChatUserEventUserID), domain.TopicUserCreated, []byte("{\"Number\":\""+fixedChatUserEventNumber.String()+"\"}"), idGen)
	parsed, err := domain.ToUserCreatedEvent(std)
	require.NoError(t, err)
	require.NotNil(t, parsed)
	require.Equal(t, fixedChatUserEventNumber, parsed.Number())

	wrong := event.NewStandardEvent(kernel.ID(fixedChatUserEventUserID), "wrong.topic", []byte("{}"), idGen)
	parsed2, err := domain.ToUserCreatedEvent(wrong)
	require.ErrorIs(t, err, errors.ErrWrongEventType)
	require.Nil(t, parsed2)
}

func TestChatUserCreatedEvent_MarshalUnmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedChatUserEventID).Times(2)
	ev, err := domain.NewUserCreatedEvent(fixedChatUserEventUserID, fixedChatUserEventNumber, idGen)
	require.NoError(t, err)
	payload, err := ev.Marshal()
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	std := event.NewStandardEvent(kernel.ID(fixedChatUserEventUserID), domain.TopicUserCreated, payload, idGen)
	converted, err := domain.ToUserCreatedEvent(std)
	require.NoError(t, err)
	require.Equal(t, ev.Number(), converted.Number())
}
