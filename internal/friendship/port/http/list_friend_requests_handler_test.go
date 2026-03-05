package http_test

import (
	"context"
	friendshipApplication "gochat/internal/friendship/application"
	friendshipDomain "gochat/internal/friendship/domain"
	friendshipDTO "gochat/internal/friendship/dto"
	friendshipHTTP "gochat/internal/friendship/port/http"
	ginMocks "gochat/internal/infrastructure/gin/mocks"
	"gochat/internal/shared/api"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"
	httptestutil "gochat/pkg/httptest"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type fakeListFriendRequestsUseCase struct {
	exec func(ctx context.Context, in *friendshipApplication.ListFriendRequestsInput) (*friendshipApplication.ListFriendRequestsOutput, error)
}

func (f fakeListFriendRequestsUseCase) Execute(ctx context.Context, in *friendshipApplication.ListFriendRequestsInput) (*friendshipApplication.ListFriendRequestsOutput, error) {
	return f.exec(ctx, in)
}

func TestListFriendRequestsRequest_Bind(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")

	tests := []struct {
		name      string
		userID    kernel.UserID
		path      string
		expectErr bool
		want      *friendshipHTTP.ListFriendRequestsRequest
	}{
		{
			name:      "success with query",
			userID:    fixedUserID,
			path:      "/friendship-requests?base_id=request-9&limit=50",
			expectErr: false,
			want: &friendshipHTTP.ListFriendRequestsRequest{
				UserID: fixedUserID,
				BaseID: kernel.OperationID("request-9"),
				Limit:  50,
			},
		},
		{
			name:      "default limit",
			userID:    fixedUserID,
			path:      "/friendship-requests",
			expectErr: false,
			want: &friendshipHTTP.ListFriendRequestsRequest{
				UserID: fixedUserID,
				BaseID: "",
				Limit:  20,
			},
		},
		{
			name:      "invalid limit",
			userID:    fixedUserID,
			path:      "/friendship-requests?limit=bad",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			require.NoError(t, err)
			req = req.WithContext(ctxutil.WithUserID(req.Context(), tt.userID))
			ginContext.Request = req

			request := &friendshipHTTP.ListFriendRequestsRequest{}
			err = request.Bind(ginContext)
			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, request)
		})
	}
}

func TestNewListFriendRequestsHandler(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")
	fixedBaseID := kernel.OperationID("request-9")
	fixedSentAt := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	fixedRequest := friendshipDomain.LoadFriendRequest(
		"request-10",
		"user-2",
		fixedUserID,
		"hello",
		friendshipDomain.StatePending,
		fixedSentAt,
	)

	tests := []struct {
		name           string
		injectUserID   kernel.UserID
		path           string
		expectValidate bool
		useCase        fakeListFriendRequestsUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps request to input and returns request list",
			injectUserID:   fixedUserID,
			path:           "/friendship-requests?base_id=request-9&limit=30",
			expectValidate: true,
			useCase: fakeListFriendRequestsUseCase{exec: func(_ context.Context, in *friendshipApplication.ListFriendRequestsInput) (*friendshipApplication.ListFriendRequestsOutput, error) {
				assert.Equal(t, fixedUserID, in.UserID)
				assert.Equal(t, fixedBaseID, in.BaseID)
				assert.Equal(t, 30, in.Limit)
				return &friendshipApplication.ListFriendRequestsOutput{Requests: []*friendshipDomain.FriendRequest{fixedRequest}}, nil
			}},
			expectResponse: api.NewResponseWithData(&friendshipHTTP.ListFriendRequestsResponseData{Requests: friendshipDTO.ToFriendRequestDTOs([]*friendshipDomain.FriendRequest{fixedRequest})}),
		},
		{
			name:           "invalid limit returns invalid param",
			injectUserID:   fixedUserID,
			path:           "/friendship-requests?limit=bad",
			expectValidate: false,
			useCase: fakeListFriendRequestsUseCase{exec: func(_ context.Context, _ *friendshipApplication.ListFriendRequestsInput) (*friendshipApplication.ListFriendRequestsOutput, error) {
				t.Fatalf("use case should not be called when bind fails")
				return nil, nil
			}},
			expectResponse: api.ResponseInvalidParam,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			router := httptestutil.NewTestRouter(t)
			router.Use(func(c *gin.Context) {
				c.Request = c.Request.WithContext(ctxutil.WithUserID(c.Request.Context(), tt.injectUserID))
				c.Next()
			})

			validator := ginMocks.NewMockValidator(ctrl)
			if tt.expectValidate {
				validator.EXPECT().Validate(gomock.Any(), gomock.Any()).Return("", nil)
			}

			router.GET("/friendship-requests", friendshipHTTP.NewListFriendRequestsHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
