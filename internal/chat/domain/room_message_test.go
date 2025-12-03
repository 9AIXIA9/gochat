package domain_test

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	eventMocks "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	kernelmocks "gochat/internal/shared/kernel/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedRoomMessageID       kernel.MessageID = "room-msg-123"
	fixedRoomMessageEvID     event.ID         = "event-room-msg-1"
	fixedRoomMessageSenderID kernel.UserID    = "room-sender"
	roomMessageContent                        = "hello room"
	roomMsgTimeTolerance                      = 150 * time.Millisecond
)

func TestRoomMessage_CreateRoomMessage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	room := domain.CreateRoom("room-x", "60001", fixedRoomMessageSenderID) // sender is owner -> member

	msgIDGen := kernelmocks.NewMockMessageIDGenerator(ctrl)
	msgIDGen.EXPECT().Generate().Return(fixedRoomMessageID)

	evIDGen := eventMocks.NewMockIDGenerator(ctrl)
	evIDGen.EXPECT().Generate().Return(fixedRoomMessageEvID)

	start := time.Now().UTC()
	rm, err := domain.CreateRoomMessage(room, fixedRoomMessageSenderID, roomMessageContent, msgIDGen, evIDGen)
	require.NoError(t, err)
	require.NotNil(t, rm)
	assert.Equal(t, fixedRoomMessageID, rm.ID())
	assert.Equal(t, fixedRoomMessageSenderID, rm.SenderID())
	assert.Equal(t, kernel.RoomID("room-x"), rm.RoomID())
	assert.Equal(t, roomMessageContent, rm.Content())
	assert.WithinDuration(t, start, rm.SentAt(), roomMsgTimeTolerance)

	events := rm.GetEvents()
	assert.Len(t, events, 1)
	assert.Empty(t, rm.GetEvents()) // drained

	createdEv := events[0]
	assert.Equal(t, fixedRoomMessageEvID, createdEv.ID())
	assert.Equal(t, domain.TopicRoomMessageCreated, createdEv.Topic())
	assert.Equal(t, kernel.ID(fixedRoomMessageID), createdEv.AggregateID())
	assert.Equal(t, []byte(""), createdEv.Payload())
}

func TestRoomMessage_CreateRoomMessage_NotMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	room := domain.CreateRoom("room-y", "60002", "room-owner") // sender not a member

	msgIDGen := kernelmocks.NewMockMessageIDGenerator(ctrl)
	// Expect no Generate call
	msgIDGen.EXPECT().Generate().Times(0)
	evIDGen := eventMocks.NewMockIDGenerator(ctrl)
	// Expect no event ID generate
	evIDGen.EXPECT().Generate().Times(0)

	rm, err := domain.CreateRoomMessage(room, fixedRoomMessageSenderID, roomMessageContent, msgIDGen, evIDGen)
	require.ErrorIs(t, err, errors.ErrNotBelongTo)
	assert.Nil(t, rm)
}

func TestRoomMessage_LoadRoomMessage(t *testing.T) {
	rm := domain.LoadRoomMessage(fixedRoomMessageID, fixedRoomMessageSenderID, "room-z", roomMessageContent, time.Now().UTC())
	require.NotNil(t, rm)
	assert.Equal(t, fixedRoomMessageID, rm.ID())
	assert.Equal(t, fixedRoomMessageSenderID, rm.SenderID())
	assert.Equal(t, kernel.RoomID("room-z"), rm.RoomID())
	assert.Equal(t, roomMessageContent, rm.Content())
}
