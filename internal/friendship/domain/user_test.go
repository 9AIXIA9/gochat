package domain_test

import (
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/domain/mocks"
	"gochat/internal/shared/event"
	eventMocks "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedFromUserID      kernel.UserID      = "from-user-123"
	fixedFromUserNumber  kernel.UserNumber  = "10086"
	fixedToUserID        kernel.UserID      = "to-user-456"
	fixedToUserNumber    kernel.UserNumber  = "10087"
	fixedFriendRequestID kernel.OperationID = "friend-req-789"
	fixedRequestContent                     = "Let's be friends!"
	fixedEventID         event.ID           = "ev-2333"
	userTimeTolerance                       = 150 * time.Millisecond
)

func TestCreateUser(t *testing.T) {
	user := domain.CreateUser(fixedFromUserID, fixedFromUserNumber)
	require.NotNil(t, user)
	assert.Equal(t, fixedFromUserID, user.ID())
	assert.Equal(t, fixedFromUserNumber, user.Number())
	assert.Empty(t, user.FriendIDs())
	assert.Empty(t, user.Requests())
	assert.Empty(t, user.GetEvents())
}

func TestUser_SendFriendRequest_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFromUser := domain.LoadUser(
		fixedFromUserID,
		fixedFromUserNumber,
		[]kernel.UserID{"mock-friend-1", "mock-friend-2"},
		[]*domain.FriendRequest{
			domain.LoadFriendRequest(
				fixedFriendRequestID,
				"mock-user-3",
				fixedFromUserID,
				fixedRequestContent,
				domain.StatePending,
				time.Now().UTC(),
			),
		},
	)

	mockOperationIDGenerator := mocks.NewMockOperationIDGenerator(ctrl)
	mockOperationIDGenerator.EXPECT().Generate().Return(fixedFriendRequestID)

	start := time.Now().UTC()

	req, err := mockFromUser.SendFriendRequest(
		fixedToUserID,
		fixedRequestContent,
		mockOperationIDGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, req)

	assert.Equal(t, fixedRequestContent, req.Content())
	assert.Equal(t, fixedToUserID, req.To())
	assert.Equal(t, fixedFromUserID, req.From())
	assert.Equal(t, domain.StatePending, req.State())
	assert.WithinDuration(t, start, req.SentAt(), userTimeTolerance)
}

func TestUser_SendFriendRequest_AlreadyFriends(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFromUser := domain.LoadUser(
		fixedFromUserID,
		fixedFromUserNumber,
		[]kernel.UserID{"mock-friend-1", "mock-friend-2", fixedToUserID},
		[]*domain.FriendRequest{
			domain.LoadFriendRequest(
				fixedFriendRequestID,
				"mock-user-3",
				fixedFromUserID,
				fixedRequestContent,
				domain.StatePending,
				time.Now().UTC(),
			),
		},
	)

	mockOperationIDGenerator := mocks.NewMockOperationIDGenerator(ctrl)
	mockOperationIDGenerator.EXPECT().Generate().Times(0)

	req, err := mockFromUser.SendFriendRequest(
		fixedToUserID,
		fixedRequestContent,
		mockOperationIDGenerator,
	)
	require.ErrorIs(t, err, domain.ErrAlreadyBeenFriends)
	require.Nil(t, req)
}

func TestUser_SendFriendRequest_AddSelfFriend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFromUser := domain.LoadUser(
		fixedFromUserID,
		fixedFromUserNumber,
		[]kernel.UserID{"mock-friend-1", "mock-friend-2"},
		[]*domain.FriendRequest{
			domain.LoadFriendRequest(
				fixedFriendRequestID,
				"mock-user-3",
				fixedFromUserID,
				fixedRequestContent,
				domain.StatePending,
				time.Now().UTC(),
			),
		},
	)

	mockOperationIDGenerator := mocks.NewMockOperationIDGenerator(ctrl)
	mockOperationIDGenerator.EXPECT().Generate().Times(0)

	req, err := mockFromUser.SendFriendRequest(
		fixedFromUserID,
		fixedRequestContent,
		mockOperationIDGenerator,
	)
	require.ErrorIs(t, err, domain.ErrAddYourselfAsFriend)
	require.Nil(t, req)
}

func TestUser_ReceiveFriendRequest_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReq := domain.LoadFriendRequest(
		fixedFriendRequestID,
		fixedFromUserID,
		fixedToUserID,
		fixedRequestContent,
		domain.StatePending,
		time.Now().UTC(),
	)

	mockToUser := domain.LoadUser(
		fixedToUserID,
		fixedToUserNumber,
		[]kernel.UserID{"mock-friend-1", "mock-friend-2"},
		[]*domain.FriendRequest{
			domain.LoadFriendRequest(
				fixedFriendRequestID,
				"mock-user-3",
				fixedToUserID,
				fixedRequestContent,
				domain.StatePending,
				time.Now().UTC(),
			),
		},
	)

	mockEventIDGenerator := eventMocks.NewMockIDGenerator(ctrl)
	mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID)

	err := mockToUser.ReceiveFriendRequest(mockReq, mockEventIDGenerator)
	require.NoError(t, err)
	assert.Contains(t, mockToUser.Requests(), mockReq)

	evs := mockToUser.GetEvents()
	assert.Len(t, evs, 1)

	receivedEv := evs[0]
	assert.Equal(t, fixedEventID, receivedEv.ID())
	assert.Equal(t, domain.TopicFriendRequestReceived, receivedEv.Topic())
	assert.Equal(t, kernel.ID(mockReq.To()), receivedEv.AggregateID())
}

func TestUser_ReceiveFriendRequest_WrongRequestTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReq := domain.LoadFriendRequest(
		fixedFriendRequestID,
		fixedFromUserID,
		"wrong-to-user",
		fixedRequestContent,
		domain.StatePending,
		time.Now().UTC(),
	)

	mockToUser := domain.LoadUser(
		fixedToUserID,
		fixedToUserNumber,
		[]kernel.UserID{"mock-friend-1", "mock-friend-2"},
		[]*domain.FriendRequest{
			domain.LoadFriendRequest(
				fixedFriendRequestID,
				"mock-user-3",
				fixedToUserID,
				fixedRequestContent,
				domain.StatePending,
				time.Now().UTC(),
			),
		},
	)

	mockEventIDGenerator := eventMocks.NewMockIDGenerator(ctrl)
	mockEventIDGenerator.EXPECT().Generate().Times(0)

	err := mockToUser.ReceiveFriendRequest(mockReq, mockEventIDGenerator)
	require.ErrorIs(t, err, domain.ErrFriendRequestNotForUser)
}

func TestUser_ReceiveFriendRequest_AlreadyFriend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReq := domain.LoadFriendRequest(
		fixedFriendRequestID,
		fixedFromUserID,
		fixedToUserID,
		fixedRequestContent,
		domain.StatePending,
		time.Now().UTC(),
	)

	mockToUser := domain.LoadUser(
		fixedToUserID,
		fixedToUserNumber,
		[]kernel.UserID{"mock-friend-1", "mock-friend-2", fixedFromUserID},
		[]*domain.FriendRequest{
			domain.LoadFriendRequest(
				fixedFriendRequestID,
				"mock-user-3",
				fixedToUserID,
				fixedRequestContent,
				domain.StatePending,
				time.Now().UTC(),
			),
		},
	)

	mockEventIDGenerator := eventMocks.NewMockIDGenerator(ctrl)
	mockEventIDGenerator.EXPECT().Generate().Times(0)

	err := mockToUser.ReceiveFriendRequest(mockReq, mockEventIDGenerator)
	require.ErrorIs(t, err, domain.ErrAlreadyBeenFriends)
}

func TestUser_ReceiveFriendRequest_AlreadySendRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReq := domain.LoadFriendRequest(
		fixedFriendRequestID,
		fixedFromUserID,
		fixedToUserID,
		fixedRequestContent,
		domain.StatePending,
		time.Now().UTC(),
	)

	mockToUser := domain.LoadUser(
		fixedToUserID,
		fixedToUserNumber,
		[]kernel.UserID{"mock-friend-1", "mock-friend-2"},
		[]*domain.FriendRequest{
			domain.LoadFriendRequest(
				fixedFriendRequestID,
				"mock-user-3",
				fixedToUserID,
				fixedRequestContent,
				domain.StatePending,
				time.Now().UTC(),
			),
			mockReq,
		},
	)

	mockEventIDGenerator := eventMocks.NewMockIDGenerator(ctrl)
	mockEventIDGenerator.EXPECT().Generate().Times(0)

	err := mockToUser.ReceiveFriendRequest(mockReq, mockEventIDGenerator)
	require.ErrorIs(t, err, domain.ErrFriendRequestExists)
}

func TestUser_AgreeFriendRequest_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockToUser := domain.LoadUser(
		fixedToUserID,
		fixedToUserNumber,
		[]kernel.UserID{"mock-friend-1", "mock-friend-2"},
		[]*domain.FriendRequest{
			domain.LoadFriendRequest(
				fixedFriendRequestID,
				fixedFromUserID,
				fixedToUserID,
				fixedRequestContent,
				domain.StatePending,
				time.Now().UTC(),
			),
		},
	)

	mockEventIDGenerator := eventMocks.NewMockIDGenerator(ctrl)
	mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1)

	err := mockToUser.AgreeFriendRequest(fixedFriendRequestID, mockEventIDGenerator)
	require.NoError(t, err)

	assert.Contains(t, mockToUser.FriendIDs(), fixedFromUserID)

	found := false
	for _, request := range mockToUser.Requests() {
		if request.ID() == fixedFriendRequestID {
			found = true
			require.Equal(t, domain.StateAgreed, request.State())
			break
		}
	}
	require.Truef(t, found, "expected friend request %s to exist", fixedFriendRequestID)

	evs := mockToUser.GetEvents()
	assert.Len(t, evs, 1)

	agreedEv := evs[0]
	assert.Equal(t, fixedEventID, agreedEv.ID())
	assert.Equal(t, domain.TopicFriendRequestAgreed, agreedEv.Topic())
	assert.Equal(t, kernel.ID(fixedFriendRequestID), agreedEv.AggregateID())
}

func TestUser_AgreeFriendRequest_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockToUser := domain.LoadUser(
		fixedToUserID,
		fixedToUserNumber,
		[]kernel.UserID{"mock-friend-1", "mock-friend-2"},
		[]*domain.FriendRequest{
			domain.LoadFriendRequest(
				"other-request-id",
				fixedFromUserID,
				fixedToUserID,
				fixedRequestContent,
				domain.StatePending,
				time.Now().UTC(),
			),
		},
	)

	mockEventIDGenerator := eventMocks.NewMockIDGenerator(ctrl)
	mockEventIDGenerator.EXPECT().Generate().Times(0)

	err := mockToUser.AgreeFriendRequest(fixedFriendRequestID, mockEventIDGenerator)
	require.NoError(t, err)

	assert.NotContains(t, mockToUser.FriendIDs(), fixedFromUserID)

	evs := mockToUser.GetEvents()
	assert.Len(t, evs, 0)
}

func TestUser_RefuseFriendRequest_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockToUser := domain.LoadUser(
		fixedToUserID,
		fixedToUserNumber,
		[]kernel.UserID{"mock-friend-1", "mock-friend-2"},
		[]*domain.FriendRequest{
			domain.LoadFriendRequest(
				fixedFriendRequestID,
				fixedFromUserID,
				fixedToUserID,
				fixedRequestContent,
				domain.StatePending,
				time.Now().UTC(),
			),
		},
	)

	mockEventIDGenerator := eventMocks.NewMockIDGenerator(ctrl)
	mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID)

	err := mockToUser.RefuseFriendRequest(fixedFriendRequestID, mockEventIDGenerator)
	require.NoError(t, err)

	assert.NotContains(t, mockToUser.FriendIDs(), fixedFromUserID)

	found := false
	for _, request := range mockToUser.Requests() {
		if request.ID() == fixedFriendRequestID {
			found = true
			require.Equal(t, domain.StateRefused, request.State())
			break
		}
	}
	require.Truef(t, found, "expected friend request %s to exist", fixedFriendRequestID)

	evs := mockToUser.GetEvents()
	assert.Len(t, evs, 1)

	refusedEv := evs[0]
	assert.Equal(t, fixedEventID, refusedEv.ID())
	assert.Equal(t, domain.TopicFriendRequestRefused, refusedEv.Topic())
	assert.Equal(t, kernel.ID(fixedFriendRequestID), refusedEv.AggregateID())
}

func TestUser_RefuseFriendRequest_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockToUser := domain.LoadUser(
		fixedToUserID,
		fixedToUserNumber,
		[]kernel.UserID{"mock-friend-1", "mock-friend-2"},
		[]*domain.FriendRequest{
			domain.LoadFriendRequest(
				"other-request-id",
				fixedFromUserID,
				fixedToUserID,
				fixedRequestContent,
				domain.StatePending,
				time.Now().UTC(),
			),
		},
	)

	mockEventIDGenerator := eventMocks.NewMockIDGenerator(ctrl)
	mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(0)

	err := mockToUser.RefuseFriendRequest(fixedFriendRequestID, mockEventIDGenerator)
	require.NoError(t, err)

	assert.NotContains(t, mockToUser.FriendIDs(), fixedFromUserID)

	evs := mockToUser.GetEvents()
	assert.Len(t, evs, 0)
}
