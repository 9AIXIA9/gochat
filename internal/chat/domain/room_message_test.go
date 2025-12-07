package domain_test

import (
	"gochat/internal/chat/domain"
	"gochat/internal/chat/domain/mocks"
	"gochat/internal/shared/kernel"
	kernelmocks "gochat/internal/shared/kernel/mocks"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const fixedRoomMemberCount = 10

func TestLoadRoomMessage(t *testing.T) {
	mockRecipients := make([]kernel.UserID, 10)
	mockStates := make(map[kernel.UserID]domain.MessageState)
	for i := 0; i < fixedRoomMemberCount; i++ {
		mockRecipients[i] = kernel.UserID("user-" + strconv.Itoa(i+1))
		mockStates[mockRecipients[i]] = domain.MessageStateRead
	}

	message := domain.LoadRoomMessage(
		fixedMessageID,
		fixedUserID,
		mockStates,
		fixedRoomID,
		"Hello, Room!",
		time.Now().UTC(),
	)

	require.NotNil(t, message)
	assert.Equal(t, fixedMessageID, message.ID())
	assert.Equal(t, fixedUserID, message.SenderID())
	assert.Equal(t, fixedRoomID, message.RoomID())
	assert.Equal(t, "Hello, Room!", message.Content())
	for id, state := range mockStates {
		assert.Equal(t, state, message.State(id))
	}
}

func TestCreateRoomMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockNotifier := mocks.NewMockRoomMessageNotifier(ctrl)

	mockRecipients := make([]kernel.UserID, 10)
	for i := 0; i < fixedRoomMemberCount; i++ {
		mockRecipients[i] = kernel.UserID("user-" + strconv.Itoa(i+1))
	}

	mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1)
	mockNotifier.EXPECT().Notify(gomock.Any(), mockRecipients).Return(mockRecipients, nil).Times(1)

	start := time.Now()
	message, err := domain.CreateRoomMessage(
		fixedRoomID,
		fixedUserID,
		mockRecipients,
		"Hello, Room!",
		mockMessageIDGenerator,
		mockNotifier,
	)

	require.NoError(t, err)
	require.NotNil(t, message)
	assert.Equal(t, fixedMessageID, message.ID())
	assert.Equal(t, fixedUserID, message.SenderID())
	assert.Equal(t, fixedRoomID, message.RoomID())
	assert.Equal(t, "Hello, Room!", message.Content())
	assert.WithinDuration(t, start, message.SentAt(), timeTolerance)

	for _, state := range message.States() {
		assert.Equal(t, domain.MessageStateDelivered, state)
	}

	// 部分未成功
	mockSuccessIDs := mockRecipients[:5]
	mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1)
	mockNotifier.EXPECT().Notify(gomock.Any(), mockRecipients).Return(mockSuccessIDs, nil).Times(1)

	message, err = domain.CreateRoomMessage(
		fixedRoomID,
		fixedUserID,
		mockRecipients,
		"Hello again, Room!",
		mockMessageIDGenerator,
		mockNotifier,
	)

	require.NoError(t, err)
	require.NotNil(t, message)

	for _, id := range mockSuccessIDs {
		assert.Equal(t, domain.MessageStateDelivered, message.State(id))
	}
	for _, id := range mockRecipients[5:] {
		assert.Equal(t, domain.MessageStateUndelivered, message.State(id))
	}
}

func TestRoomMessage_Deliver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStates := make(map[kernel.UserID]domain.MessageState, fixedRoomMemberCount+1)
	for i := 0; i < fixedRoomMemberCount; i++ {
		mockID := kernel.UserID("user-" + strconv.Itoa(i+1))
		if i < 5 {
			mockStates[mockID] = domain.MessageStateUndelivered
		} else {
			mockStates[mockID] = domain.MessageStateRead
		}
	}
	mockStates[fixedUserID] = domain.MessageStateUndelivered

	message := domain.LoadRoomMessage(
		fixedMessageID,
		fixedSenderID,
		mockStates,
		fixedRoomID,
		"Hello, Room!",
		time.Now().UTC(),
	)

	mockNotifier := mocks.NewMockRoomMessageNotifier(ctrl)
	mockNotifier.EXPECT().Notify(message, []kernel.UserID{fixedUserID}).Return([]kernel.UserID{fixedUserID}, nil).Times(1)

	err := message.Deliver(fixedUserID, mockNotifier)
	require.NoError(t, err)
	assert.Equal(t, domain.MessageStateDelivered, message.State(fixedUserID))

	// 重复投递
	err = message.Deliver(fixedUserID, mockNotifier)
	require.NoError(t, err)

	// 投递失败
	mockNotifier.EXPECT().Notify(message, []kernel.UserID{"user-1"}).Return([]kernel.UserID{}, nil).Times(1)

	err = message.Deliver("user-1", mockNotifier)
	require.NoError(t, err)
	assert.Equal(t, domain.MessageStateUndelivered, message.State("user-1"))
}
