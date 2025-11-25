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
	fixedRoomMsgCreatedEventID   event.ID         = "room-msg-ev-1"
	fixedRoomMsgCreatedMessageID kernel.MessageID = "room-msg-created-456"
)

func TestRoomMessageCreatedEvent_NewRoomMessageCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedRoomMsgCreatedEventID)

	ev, err := domain.NewRoomMessageCreatedEvent(fixedRoomMsgCreatedMessageID, idGen)
	require.NoError(t, err)
	require.NotNil(t, ev)
	require.Equal(t, fixedRoomMsgCreatedEventID, ev.ID())
	require.Equal(t, domain.TopicRoomMessageCreated, ev.Topic())
	require.Equal(t, kernel.ID(fixedRoomMsgCreatedMessageID), ev.AggregateID())
	require.Equal(t, []byte(""), ev.Payload())
}

func TestRoomMessageCreatedEvent_ToRoomMessageCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedRoomMsgCreatedEventID).AnyTimes()

	std := event.NewStandardEvent(kernel.ID(fixedRoomMsgCreatedMessageID), domain.TopicRoomMessageCreated, []byte(""), idGen)
	parsed, err := domain.ToRoomMessageCreatedEvent(std)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	wrong := event.NewStandardEvent(kernel.ID(fixedRoomMsgCreatedMessageID), "wrong.topic", []byte(""), idGen)
	parsed2, err := domain.ToRoomMessageCreatedEvent(wrong)
	require.ErrorIs(t, err, errors.ErrWrongEventType)
	require.Nil(t, parsed2)
}

func TestChatRoomMessageCreatedEvent_MarshalUnmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedRoomMsgCreatedEventID).Times(2)
	ev, err := domain.NewRoomMessageCreatedEvent(fixedRoomMsgCreatedMessageID, idGen)
	require.NoError(t, err)
	payload, err := ev.Marshal()
	require.NoError(t, err)
	require.Equal(t, 0, len(payload))
	std := event.NewStandardEvent(kernel.ID(fixedRoomMsgCreatedMessageID), domain.TopicRoomMessageCreated, payload, idGen)
	converted, err := domain.ToRoomMessageCreatedEvent(std)
	require.NoError(t, err)
	require.Equal(t, ev.ID(), converted.ID())
}
