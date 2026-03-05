package http_test

import (
	"context"
	ginMocks "gochat/internal/infrastructure/gin/mocks"
	roomshipApplication "gochat/internal/roomship/application"
	roomshipHTTP "gochat/internal/roomship/port/http"
	"gochat/internal/shared/api"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"
	httptestutil "gochat/pkg/httptest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type fakeLeaveRoomUseCase struct {
	exec func(ctx context.Context, in *roomshipApplication.LeaveRoomInput) (*kernel.NoOutput, error)
}

func (f fakeLeaveRoomUseCase) Execute(ctx context.Context, in *roomshipApplication.LeaveRoomInput) (*kernel.NoOutput, error) {
	return f.exec(ctx, in)
}

func TestLeaveRoomRequest_Bind(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")
	fixedRoomID := kernel.RoomID("room-1")

	tests := []struct {
		name   string
		userID kernel.UserID
		roomID kernel.RoomID
		want   *roomshipHTTP.LeaveRoomRequest
	}{
		{
			name:   "success",
			userID: fixedUserID,
			roomID: fixedRoomID,
			want:   &roomshipHTTP.LeaveRoomRequest{UserID: fixedUserID, RoomID: fixedRoomID},
		},
		{
			name:   "missing room id",
			userID: fixedUserID,
			roomID: "",
			want:   &roomshipHTTP.LeaveRoomRequest{UserID: fixedUserID, RoomID: ""},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)

			req, err := http.NewRequest(http.MethodDelete, "/rooms/room-1/members/me", nil)
			require.NoError(t, err)
			req = req.WithContext(ctxutil.WithUserID(req.Context(), tt.userID))
			ginContext.Request = req

			if tt.roomID != "" {
				ginContext.Params = gin.Params{{Key: "room_id", Value: string(tt.roomID)}}
			}

			request := &roomshipHTTP.LeaveRoomRequest{}
			err = request.Bind(ginContext)
			require.NoError(t, err)
			assert.Equal(t, tt.want, request)
		})
	}
}

func TestNewLeaveRoomHandler(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")
	fixedRoomID := kernel.RoomID("room-1")

	tests := []struct {
		name           string
		injectUserID   kernel.UserID
		path           string
		expectValidate bool
		useCase        fakeLeaveRoomUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps request to input",
			injectUserID:   fixedUserID,
			path:           "/rooms/room-1/members/me",
			expectValidate: true,
			useCase: fakeLeaveRoomUseCase{exec: func(_ context.Context, in *roomshipApplication.LeaveRoomInput) (*kernel.NoOutput, error) {
				assert.Equal(t, fixedUserID, in.UserID)
				assert.Equal(t, fixedRoomID, in.RoomID)
				return nil, nil
			}},
			expectResponse: api.ResponseSuccess,
		},
		{
			name:           "missing user id returns invalid param",
			injectUserID:   "",
			path:           "/rooms/room-1/members/me",
			expectValidate: true,
			useCase: fakeLeaveRoomUseCase{exec: func(_ context.Context, _ *roomshipApplication.LeaveRoomInput) (*kernel.NoOutput, error) {
				t.Fatalf("use case should not be called when input is invalid")
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

			router.DELETE("/rooms/:room_id/members/me", roomshipHTTP.NewLeaveRoomHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodDelete, tt.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
