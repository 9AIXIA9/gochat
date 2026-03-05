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

type fakeListFriendshipsUseCase struct {
	exec func(ctx context.Context, in *friendshipApplication.ListFriendshipsInput) (*friendshipApplication.ListFriendshipsOutput, error)
}

func (f fakeListFriendshipsUseCase) Execute(ctx context.Context, in *friendshipApplication.ListFriendshipsInput) (*friendshipApplication.ListFriendshipsOutput, error) {
	return f.exec(ctx, in)
}

func TestListFriendshipsRequest_Bind(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")

	tests := []struct {
		name      string
		userID    kernel.UserID
		path      string
		expectErr bool
		want      *friendshipHTTP.ListFriendshipsRequest
	}{
		{
			name:      "success with query",
			userID:    fixedUserID,
			path:      "/friendships?base_id=friendship-9&limit=50",
			expectErr: false,
			want: &friendshipHTTP.ListFriendshipsRequest{
				UserID: fixedUserID,
				BaseID: friendshipDomain.FriendshipID("friendship-9"),
				Limit:  50,
			},
		},
		{
			name:      "default limit",
			userID:    fixedUserID,
			path:      "/friendships",
			expectErr: false,
			want: &friendshipHTTP.ListFriendshipsRequest{
				UserID: fixedUserID,
				BaseID: "",
				Limit:  20,
			},
		},
		{
			name:      "invalid limit",
			userID:    fixedUserID,
			path:      "/friendships?limit=bad",
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

			request := &friendshipHTTP.ListFriendshipsRequest{}
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

func TestNewListFriendshipsHandler(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")
	fixedBaseID := friendshipDomain.FriendshipID("friendship-9")
	fixedCreatedAt := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	fixedFriendship := friendshipDomain.LoadFriendship(
		"friendship-10",
		fixedUserID,
		"user-2",
		fixedCreatedAt,
	)

	tests := []struct {
		name           string
		injectUserID   kernel.UserID
		path           string
		expectValidate bool
		useCase        fakeListFriendshipsUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps request to input and returns friendship list",
			injectUserID:   fixedUserID,
			path:           "/friendships?base_id=friendship-9&limit=30",
			expectValidate: true,
			useCase: fakeListFriendshipsUseCase{exec: func(_ context.Context, in *friendshipApplication.ListFriendshipsInput) (*friendshipApplication.ListFriendshipsOutput, error) {
				assert.Equal(t, fixedUserID, in.UserID)
				assert.Equal(t, fixedBaseID, in.BaseID)
				assert.Equal(t, 30, in.Limit)
				return &friendshipApplication.ListFriendshipsOutput{Friendships: []*friendshipDomain.Friendship{fixedFriendship}}, nil
			}},
			expectResponse: api.NewResponseWithData(&friendshipHTTP.ListFriendshipsResponseData{Friendships: friendshipDTO.ToFriendshipDTOs([]*friendshipDomain.Friendship{fixedFriendship})}),
		},
		{
			name:           "invalid limit returns invalid param",
			injectUserID:   fixedUserID,
			path:           "/friendships?limit=bad",
			expectValidate: false,
			useCase: fakeListFriendshipsUseCase{exec: func(_ context.Context, _ *friendshipApplication.ListFriendshipsInput) (*friendshipApplication.ListFriendshipsOutput, error) {
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

			router.GET("/friendships", friendshipHTTP.NewListFriendshipsHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
