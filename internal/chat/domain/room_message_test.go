package domain_test

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/domain/mocks"
	eventMock "gochat/internal/shared/event/mocks"
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
	mockRecipients := make([]kernel.UserID, fixedRoomMemberCount)
	for i := 0; i < fixedRoomMemberCount; i++ {
		mockRecipients[i] = kernel.UserID("user-" + strconv.Itoa(i+1))
	}

	message := domain.LoadRoomMessage(
		fixedMessageID,
		fixedUserID,
		mockRecipients,
		fixedRoomID,
		"Hello, Room!",
		time.Now().UTC(),
	)

	require.NotNil(t, message)
	assert.Equal(t, fixedMessageID, message.ID())
	assert.Equal(t, fixedUserID, message.SenderID())
	assert.Equal(t, fixedRoomID, message.RoomID())
	assert.Equal(t, "Hello, Room!", message.Content())
	assert.Equal(t, mockRecipients, message.RecipientIDs())
}

func TestCreateRoomMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockRoomshipsFinderByRoomID(ctrl)
	mockMessageIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockEventIDGenerator := eventMock.NewMockIDGenerator(ctrl)

	roomships := make([]*domain.Roomship, 0, fixedRoomMemberCount+1)
	mockRecipients := make([]kernel.UserID, 0, fixedRoomMemberCount)
	for i := 0; i < fixedRoomMemberCount; i++ {
		roomships = append(roomships, domain.LoadRoomship(
			domain.RoomshipID("roomship-"+strconv.Itoa(i+1)),
			kernel.UserID("user-"+strconv.Itoa(i+1)),
			fixedRoomID,
		))
		mockRecipients = append(mockRecipients, roomships[i].UserID())
	}
	roomships = append(roomships, domain.LoadRoomship(
		fixedRoomshipID,
		fixedUserID,
		fixedRoomID,
	))
	roomshipsWithoutSender := roomships[:fixedRoomMemberCount]

	gomock.InOrder(
		// 正常情况
		finder.EXPECT().FindsByRoomID(gomock.Any(), fixedRoomID).Return(roomships, nil).Times(1),
		mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1),
		mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1),

		// content 为空 什么都不做直接错误

		// 不是成员
		finder.EXPECT().FindsByRoomID(gomock.Any(), fixedRoomID).Return(roomshipsWithoutSender, nil).Times(1),

		// 只有自己
		finder.EXPECT().FindsByRoomID(gomock.Any(), fixedRoomID).Return([]*domain.Roomship{domain.LoadRoomship(
			fixedRoomshipID,
			fixedUserID,
			fixedRoomID,
		)}, nil).Times(1),
		mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1),
		mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1),

		// room not found
		finder.EXPECT().FindsByRoomID(gomock.Any(), fixedRoomID).Return([]*domain.Roomship{}, nil).Times(1),
	)

	// 正常情况
	start := time.Now().UTC()
	message, err := domain.CreateRoomMessage(
		context.Background(),
		fixedRoomID,
		fixedUserID,
		"Hello, Room!",
		finder,
		mockMessageIDGenerator,
		mockEventIDGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, message)
	assert.Equal(t, fixedMessageID, message.ID())
	assert.Equal(t, fixedUserID, message.SenderID())
	assert.Equal(t, fixedRoomID, message.RoomID())
	assert.Equal(t, "Hello, Room!", message.Content())
	assert.WithinDuration(t, start, message.SentAt(), timeTolerance)
	assert.Equal(t, mockRecipients, message.RecipientIDs())
	evs := message.GetEvents()
	require.Len(t, evs, 1)
	assert.Equal(t, domain.TopicRoomMessageCreated, evs[0].Topic())
	assert.Equal(t, fixedEventID, evs[0].ID())

	// content 为空
	message, err = domain.CreateRoomMessage(
		context.Background(),
		fixedRoomID,
		fixedUserID,
		"",
		finder,
		mockMessageIDGenerator,
		mockEventIDGenerator,
	)
	require.ErrorIs(t, err, domain.ErrEmptyMessageContent)
	require.Nil(t, message)

	// 不是成员
	_, err = domain.CreateRoomMessage(
		context.Background(),
		fixedRoomID,
		fixedUserID,
		"Hello, Room!",
		finder,
		mockMessageIDGenerator,
		mockEventIDGenerator,
	)
	require.ErrorIs(t, err, domain.ErrNotMember)

	// 只有自己
	start = time.Now().UTC()
	message, err = domain.CreateRoomMessage(
		context.Background(),
		fixedRoomID,
		fixedUserID,
		"Hello to myself!",
		finder,
		mockMessageIDGenerator,
		mockEventIDGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, message)
	assert.Equal(t, fixedMessageID, message.ID())
	assert.Equal(t, fixedUserID, message.SenderID())
	assert.Equal(t, fixedRoomID, message.RoomID())
	assert.Equal(t, "Hello to myself!", message.Content())
	assert.WithinDuration(t, start, message.SentAt(), timeTolerance)
	assert.Empty(t, message.RecipientIDs())
	evs = message.GetEvents()
	require.Len(t, evs, 1)
	assert.Equal(t, domain.TopicRoomMessageCreated, evs[0].Topic())
	assert.Equal(t, fixedEventID, evs[0].ID())

	// room not found
	message, err = domain.CreateRoomMessage(
		context.Background(),
		fixedRoomID,
		fixedUserID,
		"Hello!",
		finder,
		mockMessageIDGenerator,
		mockEventIDGenerator,
	)
	require.ErrorIs(t, err, domain.ErrRoomNotFound)
	require.Nil(t, message)
}
