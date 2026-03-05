package http_test

import (
	"context"
	ginMocks "gochat/internal/infrastructure/gin/mocks"
	roomshipApplication "gochat/internal/roomship/application"
	roomshipDomain "gochat/internal/roomship/domain"
	roomshipDTO "gochat/internal/roomship/dto"
	roomshipHTTP "gochat/internal/roomship/port/http"
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

type fakeListMemberRequestsUseCase struct {
	exec func(ctx context.Context, in *roomshipApplication.ListMemberRequestsInput) (*roomshipApplication.ListMemberRequestsOutput, error)
}

func (f fakeListMemberRequestsUseCase) Execute(ctx context.Context, in *roomshipApplication.ListMemberRequestsInput) (*roomshipApplication.ListMemberRequestsOutput, error) {
	return f.exec(ctx, in)
}

func TestListMemberRequestsRequest_Bind(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")

	tests := []struct {
		name      string
		userID    kernel.UserID
		path      string
		expectErr bool
		want      *roomshipHTTP.ListMemberRequestsRequest
	}{
		{
			name:      "success with query",
			userID:    fixedUserID,
			path:      "/rooms/requests?base_id=request-9&limit=50",
			expectErr: false,
			want:      &roomshipHTTP.ListMemberRequestsRequest{UserID: fixedUserID, BaseID: kernel.OperationID("request-9"), Limit: 50},
		},
		{
			name:      "default limit",
			userID:    fixedUserID,
			path:      "/rooms/requests",
			expectErr: false,
			want:      &roomshipHTTP.ListMemberRequestsRequest{UserID: fixedUserID, BaseID: "", Limit: 20},
		},
		{
			name:      "invalid limit",
			userID:    fixedUserID,
			path:      "/rooms/requests?limit=bad",
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

			request := &roomshipHTTP.ListMemberRequestsRequest{}
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

func TestNewListMemberRequestsHandler(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")
	fixedBaseID := kernel.OperationID("request-9")
	fixedTime := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	fixedRequest := roomshipDomain.LoadMemberRequest(
		"request-10",
		roomshipDomain.StatePending,
		"user-2",
		"room-1",
		"join please",
		"",
		fixedTime,
		fixedTime,
	)

	tests := []struct {
		name           string
		injectUserID   kernel.UserID
		path           string
		expectValidate bool
		useCase        fakeListMemberRequestsUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps request to input and returns member request list",
			injectUserID:   fixedUserID,
			path:           "/rooms/requests?base_id=request-9&limit=30",
			expectValidate: true,
			useCase: fakeListMemberRequestsUseCase{exec: func(_ context.Context, in *roomshipApplication.ListMemberRequestsInput) (*roomshipApplication.ListMemberRequestsOutput, error) {
				assert.Equal(t, fixedUserID, in.UserID)
				assert.Equal(t, fixedBaseID, in.BaseID)
				assert.Equal(t, 30, in.Limit)
				return &roomshipApplication.ListMemberRequestsOutput{Requests: []*roomshipDomain.MemberRequest{fixedRequest}}, nil
			}},
			expectResponse: api.NewResponseWithData(&roomshipHTTP.ListMemberRequestsResponseData{Requests: roomshipDTO.ToMemberRequestDTOs([]*roomshipDomain.MemberRequest{fixedRequest})}),
		},
		{
			name:           "invalid limit returns invalid param",
			injectUserID:   fixedUserID,
			path:           "/rooms/requests?limit=bad",
			expectValidate: false,
			useCase: fakeListMemberRequestsUseCase{exec: func(_ context.Context, _ *roomshipApplication.ListMemberRequestsInput) (*roomshipApplication.ListMemberRequestsOutput, error) {
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

			router.GET("/rooms/requests", roomshipHTTP.NewListMemberRequestsHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
