package http_test

import (
	"bytes"
	"context"
	chatApplication "gochat/internal/chat/application"
	chatHTTP "gochat/internal/chat/port/http"
	ginMocks "gochat/internal/infrastructure/gin/mocks"
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

type fakeReadRoomMessagesUseCase struct {
	exec func(ctx context.Context, in *chatApplication.ReadRoomMessagesInput) (*kernel.NoOutput, error)
}

func (f fakeReadRoomMessagesUseCase) Execute(ctx context.Context, in *chatApplication.ReadRoomMessagesInput) (*kernel.NoOutput, error) {
	return f.exec(ctx, in)
}

func TestReadRoomMessagesRequest_Bind(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")

	tests := []struct {
		name      string
		userID    kernel.UserID
		body      string
		expectErr bool
		want      *chatHTTP.ReadRoomMessagesRequest
	}{
		{
			name:      "success",
			userID:    fixedUserID,
			body:      `{"room_id":"room-1"}`,
			expectErr: false,
			want: &chatHTTP.ReadRoomMessagesRequest{
				UserID: fixedUserID,
				RoomID: kernel.RoomID("room-1"),
			},
		},
		{
			name:      "malformed json",
			userID:    fixedUserID,
			body:      `{"room_id":`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)

			req, err := http.NewRequest(http.MethodPut, "/chats/rooms/messages/read", bytes.NewBufferString(tt.body))
			require.NoError(t, err)
			req = req.WithContext(ctxutil.WithUserID(req.Context(), tt.userID))
			req.Header.Set("Content-Type", "application/json")
			ginContext.Request = req

			request := &chatHTTP.ReadRoomMessagesRequest{}
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

func TestNewReadRoomMessagesHandler(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")
	body := `{"room_id":"room-1"}`

	tests := []struct {
		name           string
		injectUserID   kernel.UserID
		body           string
		expectValidate bool
		useCase        fakeReadRoomMessagesUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps request to input",
			injectUserID:   fixedUserID,
			body:           body,
			expectValidate: true,
			useCase: fakeReadRoomMessagesUseCase{exec: func(_ context.Context, in *chatApplication.ReadRoomMessagesInput) (*kernel.NoOutput, error) {
				assert.Equal(t, fixedUserID, in.UserID)
				assert.Equal(t, kernel.RoomID("room-1"), in.RoomID)
				return nil, nil
			}},
			expectResponse: api.ResponseSuccess,
		},
		{
			name:           "malformed json returns invalid param",
			injectUserID:   fixedUserID,
			body:           `{"room_id":`,
			expectValidate: false,
			useCase: fakeReadRoomMessagesUseCase{exec: func(_ context.Context, _ *chatApplication.ReadRoomMessagesInput) (*kernel.NoOutput, error) {
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

			router.PUT("/chats/rooms/messages/read", chatHTTP.NewReadRoomMessagesHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodPut, "/chats/rooms/messages/read", bytes.NewBufferString(tt.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
