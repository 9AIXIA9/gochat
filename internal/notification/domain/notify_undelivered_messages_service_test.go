package domain_test

import (
	"testing"
	"time"

	"gochat/internal/notification/domain"
	"gochat/internal/notification/domain/mocks"
	"gochat/internal/shared/kernel"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	userIDNotify kernel.UserID = "user-notify"
)

func TestNotifyUndeliveredMessages_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pm := domain.ReceivePrivateMessage("pm-1", "sender-1", userIDNotify, "hi", time.Now().UTC())
	rm := domain.ReceiveRoomMessage("rm-1", "sender-2", "room-1", []kernel.UserID{userIDNotify}, "hi room", time.Now().UTC())

	pmNotifier := mocks.NewMockPrivateMessageNotifier(ctrl)
	rmNotifier := mocks.NewMockRoomMessageNotifier(ctrl)
	pmNotifier.EXPECT().NotifyPrivateMessage(pm).Return(nil)
	rmNotifier.EXPECT().NotifyRoomMessage(rm, []kernel.UserID{userIDNotify}).Return([]kernel.UserID{userIDNotify}, nil)

	err := domain.NotifyUndeliveredMessages(userIDNotify, []*domain.PrivateMessage{pm}, []*domain.RoomMessage{rm}, pmNotifier, rmNotifier)
	require.NoError(t, err)
	assert.Equal(t, domain.MessageStateDelivered, pm.State())
	assert.Equal(t, domain.MessageStateDelivered, rm.States()[userIDNotify])
}

func TestNotifyUndeliveredMessages_PrivateMessageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	pm := domain.ReceivePrivateMessage("pm-1", "sender-1", userIDNotify, "hi", time.Now().UTC())
	pmNotifier := mocks.NewMockPrivateMessageNotifier(ctrl)
	rmNotifier := mocks.NewMockRoomMessageNotifier(ctrl)
	pmNotifier.EXPECT().NotifyPrivateMessage(pm).Return(assert.AnError)
	err := domain.NotifyUndeliveredMessages(userIDNotify, []*domain.PrivateMessage{pm}, nil, pmNotifier, rmNotifier)
	require.Error(t, err)
	assert.Equal(t, domain.MessageStateUndelivered, pm.State())
}

func TestNotifyUndeliveredMessages_RoomMessageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	rm := domain.ReceiveRoomMessage("rm-1", "sender-1", "room-1", []kernel.UserID{userIDNotify}, "hi room", time.Now().UTC())
	pmNotifier := mocks.NewMockPrivateMessageNotifier(ctrl)
	rmNotifier := mocks.NewMockRoomMessageNotifier(ctrl)
	rmNotifier.EXPECT().NotifyRoomMessage(rm, []kernel.UserID{userIDNotify}).Return(nil, assert.AnError)
	err := domain.NotifyUndeliveredMessages(userIDNotify, nil, []*domain.RoomMessage{rm}, pmNotifier, rmNotifier)
	require.Error(t, err)
	assert.Equal(t, domain.MessageStateUndelivered, rm.States()[userIDNotify])
}
