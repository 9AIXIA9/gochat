package domain_test

import (
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/domain/mocks"
	eventMock "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoadFriendRequest(t *testing.T) {
	start := time.Now()
	req := domain.LoadFriendRequest(
		fixedRequestID,
		fixedUserID,
		fixedToUserID,
		fixedContent,
		domain.StatePending,
		time.Now().UTC(),
	)
	require.NotNil(t, req)
	assert.Equal(t, fixedRequestID, req.ID())
	assert.Equal(t, fixedUserID, req.From())
	assert.Equal(t, fixedToUserID, req.To())
	assert.Equal(t, fixedContent, req.Content())
	assert.Equal(t, domain.StatePending, req.State())
	assert.WithinDuration(t, start, req.SentAt(), timeTolerance)
}

func TestCreateFriendRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//正常情况
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1)

	mockOperationIDGenerator := mocks.NewMockOperationIDGenerator(ctrl)
	mockOperationIDGenerator.EXPECT().Generate().Return(fixedRequestID).Times(1)

	start := time.Now()

	req, err := domain.CreateFriendRequest(fixedUserID, fixedToUserID, fixedContent, mockOperationIDGenerator, mockIDGenerator)
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, fixedRequestID, req.ID())
	assert.Equal(t, fixedUserID, req.From())
	assert.Equal(t, fixedToUserID, req.To())
	assert.Equal(t, fixedContent, req.Content())
	assert.Equal(t, domain.StatePending, req.State())
	assert.WithinDuration(t, start, req.SentAt(), timeTolerance)

	evs := req.GetEvents()
	require.Len(t, evs, 1)

	createdEv := evs[0]
	assert.Equal(t, fixedEventID, createdEv.ID())
	assert.Equal(t, domain.TopicFriendRequestCreated, createdEv.Topic())
	assert.Equal(t, kernel.ID(fixedRequestID), createdEv.AggregateID())

	//content太长
	tooLongContent := ""
	for i := 0; i < maxContentLength+1; i++ {
		tooLongContent += "a"
	}

	req2, err2 := domain.CreateFriendRequest(fixedUserID, fixedToUserID, tooLongContent, mockOperationIDGenerator, mockIDGenerator)
	require.ErrorIs(t, err2, domain.ErrFriendRequestContentTooLong)
	require.Nil(t, req2)

	//添加自己为好友
	req3, err3 := domain.CreateFriendRequest(fixedUserID, fixedUserID, fixedContent, mockOperationIDGenerator, mockIDGenerator)
	require.ErrorIs(t, err3, domain.ErrAddYourselfAsFriend)
	require.Nil(t, req3)
}

func TestFriendRequest_Agree(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//正常情况
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1)

	mockRequest := domain.LoadFriendRequest(
		fixedRequestID,
		fixedUserID,
		fixedToUserID,
		fixedContent,
		domain.StatePending,
		time.Now().UTC(),
	)

	err := mockRequest.Agree(mockIDGenerator)
	require.NoError(t, err)
	require.Equal(t, domain.StateAgreed, mockRequest.State())

	evs := mockRequest.GetEvents()
	require.Len(t, evs, 1)

	agreedEv := evs[0]
	assert.Equal(t, fixedEventID, agreedEv.ID())
	assert.Equal(t, domain.TopicFriendRequestAgreed, agreedEv.Topic())
	assert.Equal(t, kernel.ID(fixedRequestID), agreedEv.AggregateID())

	//重复同意
	err2 := mockRequest.Agree(mockIDGenerator)
	require.ErrorIs(t, err2, domain.ErrFriendRequestHasBeenHandled)
	require.Len(t, mockRequest.GetEvents(), 0)

	//失败后又同意
	mockRequestRefused := domain.LoadFriendRequest(
		fixedRequestID,
		fixedUserID,
		fixedToUserID,
		fixedContent,
		domain.StateRefused,
		time.Now().UTC(),
	)
	err3 := mockRequestRefused.Agree(mockIDGenerator)
	require.ErrorIs(t, err3, domain.ErrFriendRequestHasBeenHandled)
	require.Len(t, mockRequestRefused.GetEvents(), 0)
}

func TestFriendRequest_Refuse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//正常情况
	mockRequest := domain.LoadFriendRequest(
		fixedRequestID,
		fixedUserID,
		fixedToUserID,
		fixedContent,
		domain.StatePending,
		time.Now().UTC(),
	)

	err := mockRequest.Refuse()
	require.NoError(t, err)
	require.Equal(t, domain.StateRefused, mockRequest.State())

	evs := mockRequest.GetEvents()
	require.Len(t, evs, 0)

	//重复同意
	err2 := mockRequest.Refuse()
	require.ErrorIs(t, err2, domain.ErrFriendRequestHasBeenHandled)
	require.Len(t, mockRequest.GetEvents(), 0)

	//失败后又同意
	mockRequestRefused := domain.LoadFriendRequest(
		fixedRequestID,
		fixedUserID,
		fixedToUserID,
		fixedContent,
		domain.StateRefused,
		time.Now().UTC(),
	)
	err3 := mockRequestRefused.Refuse()
	require.ErrorIs(t, err3, domain.ErrFriendRequestHasBeenHandled)
	require.Len(t, mockRequestRefused.GetEvents(), 0)
}
