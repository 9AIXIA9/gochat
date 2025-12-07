package domain_test

import (
	"gochat/internal/roomship/domain"
	eventMock "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	kernelmocks "gochat/internal/shared/kernel/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoadMemberRequest(t *testing.T) {
	start := time.Now()
	req := domain.LoadMemberRequest(
		fixedOperationID,
		domain.StateAgreed,
		fixedUserID,
		fixedRoomID,
		fixedContent,
		fixedOperatorID,
		time.Now().UTC(),
		time.Now().UTC(),
	)
	require.NotNil(t, req)
	assert.WithinDuration(t, start, req.CreatedAt(), timeTolerance)
	assert.WithinDuration(t, start, req.OperatedAt(), timeTolerance)
	assert.Equal(t, fixedOperationID, req.ID())
	assert.Equal(t, fixedUserID, req.ApplicantID())
	assert.Equal(t, fixedRoomID, req.RoomID())
	assert.Equal(t, fixedContent, req.Content())
	assert.Equal(t, domain.StateAgreed, req.State())
	assert.Equal(t, fixedOperatorID, req.OperatorID())
}

func TestCreateMemberRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//正常情况
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1)

	mockOperationIDGenerator := kernelmocks.NewMockOperationIDGenerator(ctrl)
	mockOperationIDGenerator.EXPECT().Generate().Return(fixedOperationID).Times(1)

	start := time.Now()

	req, err := domain.CreateMemberRequest(
		fixedUserID,
		fixedRoomID,
		fixedContent,
		mockOperationIDGenerator,
		mockIDGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, fixedOperationID, req.ID())
	assert.Equal(t, fixedUserID, req.ApplicantID())
	assert.Equal(t, fixedRoomID, req.RoomID())
	assert.Equal(t, fixedContent, req.Content())
	assert.Equal(t, domain.StatePending, req.State())
	assert.WithinDuration(t, start, req.CreatedAt(), timeTolerance)

	evs := req.GetEvents()
	require.Len(t, evs, 1)

	createdEv := evs[0]
	assert.Equal(t, fixedEventID, createdEv.ID())
	assert.Equal(t, domain.TopicMemberRequestCreated, createdEv.Topic())
	assert.Equal(t, kernel.ID(fixedOperationID), createdEv.AggregateID())

	//content太长
	tooLongContent := ""
	for i := 0; i < maxContentLength+1; i++ {
		tooLongContent += "a"
	}

	req2, err2 := domain.CreateMemberRequest(fixedUserID, fixedRoomID, tooLongContent, mockOperationIDGenerator, mockIDGenerator)
	require.ErrorIs(t, err2, domain.ErrContentTooLong)
	require.Nil(t, req2)
}

func TestMemberRequest_Agree(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//正常情况
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1)

	mockRequest := domain.LoadMemberRequest(
		fixedOperationID,
		domain.StatePending,
		fixedUserID,
		fixedRoomID,
		fixedContent,
		"",
		time.Time{},
		time.Now().UTC(),
	)

	err := mockRequest.Agree(fixedOperatorID, mockIDGenerator)
	require.NoError(t, err)
	require.Equal(t, domain.StateAgreed, mockRequest.State())
	assert.Equal(t, fixedOperatorID, mockRequest.OperatorID())
	assert.WithinDuration(t, time.Now().UTC(), mockRequest.OperatedAt(), timeTolerance)

	evs := mockRequest.GetEvents()
	require.Len(t, evs, 1)

	agreedEv := evs[0]
	assert.Equal(t, fixedEventID, agreedEv.ID())
	assert.Equal(t, domain.TopicMemberRequestAgreed, agreedEv.Topic())
	assert.Equal(t, kernel.ID(fixedOperationID), agreedEv.AggregateID())

	//重复操作
	err2 := mockRequest.Agree(fixedOperatorID, mockIDGenerator)
	require.ErrorIs(t, err2, domain.ErrHandleNotPendingRequest)
	require.Len(t, mockRequest.GetEvents(), 0)
}

func TestMemberRequest_Refuse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//正常情况
	mockRequest := domain.LoadMemberRequest(
		fixedOperationID,
		domain.StatePending,
		fixedUserID,
		fixedRoomID,
		fixedContent,
		"",
		time.Time{},
		time.Now().UTC(),
	)

	err := mockRequest.Refuse(fixedOperatorID)
	require.NoError(t, err)
	require.Equal(t, domain.StateRefused, mockRequest.State())
	assert.Len(t, mockRequest.GetEvents(), 0)
	assert.Equal(t, fixedOperatorID, mockRequest.OperatorID())
	assert.WithinDuration(t, time.Now().UTC(), mockRequest.OperatedAt(), timeTolerance)

	err2 := mockRequest.Refuse(fixedOperatorID)
	require.ErrorIs(t, err2, domain.ErrHandleNotPendingRequest)
}
