package http_test

import (
	"bytes"
	"context"
	ginMocks "gochat/internal/infrastructure/gin/mocks"
	profileApplication "gochat/internal/profile/application"
	profileHTTP "gochat/internal/profile/port/http"
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

type fakeUpdateRoomProfileUseCase struct {
	exec func(ctx context.Context, in *profileApplication.UpdateRoomProfileInput) (*kernel.NoOutput, error)
}

func (f fakeUpdateRoomProfileUseCase) Execute(ctx context.Context, in *profileApplication.UpdateRoomProfileInput) (*kernel.NoOutput, error) {
	return f.exec(ctx, in)
}

func TestUpdateRoomProfileRequest_Bind(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")

	tests := []struct {
		name      string
		userID    kernel.UserID
		params    gin.Params
		body      string
		expectErr bool
		want      *profileHTTP.UpdateRoomProfileRequest
	}{
		{
			name:   "success",
			userID: fixedUserID,
			params: gin.Params{{Key: "room_id", Value: "room-1"}},
			body:   `{"name":"family","introduction":"our room"}`,
			want: &profileHTTP.UpdateRoomProfileRequest{
				UserID: fixedUserID,
				RoomID: kernel.RoomID("room-1"),
				Body: profileHTTP.UpdateRoomProfileRequestBody{
					Name:         "family",
					Introduction: "our room",
				},
			},
		},
		{
			name:      "malformed json",
			userID:    fixedUserID,
			params:    gin.Params{{Key: "room_id", Value: "room-1"}},
			body:      `{"name":`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Params = tt.params

			req, err := http.NewRequest(http.MethodPut, "/profiles/rooms/room-1", bytes.NewBufferString(tt.body))
			require.NoError(t, err)
			req = req.WithContext(ctxutil.WithUserID(req.Context(), tt.userID))
			req.Header.Set("Content-Type", "application/json")
			ginContext.Request = req

			request := &profileHTTP.UpdateRoomProfileRequest{}
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

func TestNewUpdateRoomProfileHandler(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")
	body := `{"name":"family","introduction":"our room"}`

	tests := []struct {
		name           string
		injectUserID   kernel.UserID
		path           string
		body           string
		expectValidate bool
		useCase        fakeUpdateRoomProfileUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps context and uri and body",
			injectUserID:   fixedUserID,
			path:           "/profiles/rooms/room-1",
			body:           body,
			expectValidate: true,
			useCase: fakeUpdateRoomProfileUseCase{exec: func(_ context.Context, in *profileApplication.UpdateRoomProfileInput) (*kernel.NoOutput, error) {
				assert.Equal(t, fixedUserID, in.UserID)
				assert.Equal(t, kernel.RoomID("room-1"), in.RoomID)
				assert.Equal(t, "family", in.Name)
				assert.Equal(t, "our room", in.Introduction)
				return nil, nil
			}},
			expectResponse: api.ResponseSuccess,
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

			router.PUT("/profiles/rooms/:room_id", profileHTTP.NewUpdateRoomProfileHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodPut, tt.path, bytes.NewBufferString(tt.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
