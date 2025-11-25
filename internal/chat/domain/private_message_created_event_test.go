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
	fixedPrivateMsgCreatedEventID   event.ID         = "priv-msg-ev-1"
	fixedPrivateMsgCreatedMessageID kernel.MessageID = "priv-msg-created-123"
)

func TestPrivateMessageCreatedEvent_NewPrivateMessageCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedPrivateMsgCreatedEventID)

	ev, err := domain.NewPrivateMessageCreatedEvent(fixedPrivateMsgCreatedMessageID, idGen)
	require.NoError(t, err)
	require.NotNil(t, ev)
	require.Equal(t, fixedPrivateMsgCreatedEventID, ev.ID())
	require.Equal(t, domain.TopicPrivateMessageCreated, ev.Topic())
	require.Equal(t, kernel.ID(fixedPrivateMsgCreatedMessageID), ev.AggregateID())
	require.Equal(t, []byte(""), ev.Payload())
}

func TestPrivateMessageCreatedEvent_ToPrivateMessageCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedPrivateMsgCreatedEventID).AnyTimes()

	std := event.NewStandardEvent(kernel.ID(fixedPrivateMsgCreatedMessageID), domain.TopicPrivateMessageCreated, []byte(""), idGen)
	parsed, err := domain.ToPrivateMessageCreatedEvent(std)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	wrong := event.NewStandardEvent(kernel.ID(fixedPrivateMsgCreatedMessageID), "wrong.topic", []byte(""), idGen)
	parsed2, err := domain.ToPrivateMessageCreatedEvent(wrong)
	require.ErrorIs(t, err, errors.ErrWrongEventType)
	require.Nil(t, parsed2)
}

func TestPrivateMessageCreatedEvent_MarshalUnmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedPrivateMsgCreatedEventID).Times(2)
	ev, err := domain.NewPrivateMessageCreatedEvent(fixedPrivateMsgCreatedMessageID, idGen)
	require.NoError(t, err)
	payload, err := ev.Marshal()
	require.NoError(t, err)
	require.Equal(t, 0, len(payload))
	std := event.NewStandardEvent(kernel.ID(fixedPrivateMsgCreatedMessageID), domain.TopicPrivateMessageCreated, payload, idGen)
	converted, err := domain.ToPrivateMessageCreatedEvent(std)
	require.NoError(t, err)
	require.Equal(t, ev.ID(), converted.ID())
}
