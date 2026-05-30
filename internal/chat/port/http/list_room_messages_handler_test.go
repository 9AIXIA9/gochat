package http_test

import (
	"context"
	chatApplication "gochat/internal/chat/application"
	chatDomain "gochat/internal/chat/domain"
	chatDTO "gochat/internal/chat/dto"
	chatHTTP "gochat/internal/chat/port/http"
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

type fakeListRoomMessagesUseCase struct {
	exec func(ctx context.Context, in *chatApplication.ListRoomMessagesInput) (*chatApplication.ListRoomMessagesOutput, error)
}

func (f fakeListRoomMessagesUseCase) Execute(ctx context.Context, in *chatApplication.ListRoomMessagesInput) (*chatApplication.ListRoomMessagesOutput, error) {
	return f.exec(ctx, in)
}

func TestListRoomMessagesRequest_Bind(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")

	tests := []struct {
		name      string
		userID    kernel.UserID
		path      string
		params    gin.Params
		expectErr bool
		want      *chatHTTP.ListRoomMessagesRequest
	}{
		{
			name:      "success with query and uri",
			userID:    fixedUserID,
			path:      "/chats/rooms/messages/room-1?base_id=message-9&limit=30",
			params:    gin.Params{{Key: "room_id", Value: "room-1"}},
			expectErr: false,
			want: &chatHTTP.ListRoomMessagesRequest{
				UserID: fixedUserID,
				RoomID: kernel.RoomID("room-1"),
				BaseID: kernel.MessageID("message-9"),
				Limit:  30,
			},
		},
		{
			name:      "default limit",
			userID:    fixedUserID,
			path:      "/chats/rooms/messages/room-1",
			params:    gin.Params{{Key: "room_id", Value: "room-1"}},
			expectErr: false,
			want: &chatHTTP.ListRoomMessagesRequest{
				UserID: fixedUserID,
				RoomID: kernel.RoomID("room-1"),
				BaseID: "",
				Limit:  20,
			},
		},
		{
			name:      "invalid limit",
			userID:    fixedUserID,
			path:      "/chats/rooms/messages/room-1?limit=bad",
			params:    gin.Params{{Key: "room_id", Value: "room-1"}},
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
			ginContext.Params = tt.params

			request := &chatHTTP.ListRoomMessagesRequest{}
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

func TestNewListRoomMessagesHandler(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")
	fixedRoomID := kernel.RoomID("room-1")
	fixedBaseID := kernel.MessageID("message-9")
	fixedSentAt := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	fixedRoomMessage := chatDomain.LoadRoomMessage(
		"message-10",
		fixedUserID,
		[]kernel.UserID{"user-2"},
		fixedRoomID,
		"hello room",
		fixedSentAt,
	)

	tests := []struct {
		name           string
		injectUserID   kernel.UserID
		path           string
		expectValidate bool
		useCase        fakeListRoomMessagesUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps request to input and returns room messages",
			injectUserID:   fixedUserID,
			path:           "/chats/rooms/messages/room-1?base_id=message-9&limit=30",
			expectValidate: true,
			useCase: fakeListRoomMessagesUseCase{exec: func(_ context.Context, in *chatApplication.ListRoomMessagesInput) (*chatApplication.ListRoomMessagesOutput, error) {
				assert.Equal(t, fixedUserID, in.UserID)
				assert.Equal(t, fixedRoomID, in.RoomID)
				assert.Equal(t, fixedBaseID, in.BaseID)
				assert.Equal(t, 30, in.Limit)
				return &chatApplication.ListRoomMessagesOutput{RoomMessages: []*chatDomain.RoomMessage{fixedRoomMessage}}, nil
			}},
			expectResponse: api.NewResponseWithData(&chatHTTP.ListRoomMessagesResponseData{RoomMessages: chatDTO.ToRoomMessageDTOs([]*chatDomain.RoomMessage{fixedRoomMessage})}),
		},
		{
			name:           "invalid query returns invalid param",
			injectUserID:   fixedUserID,
			path:           "/chats/rooms/messages/room-1?limit=bad",
			expectValidate: false,
			useCase: fakeListRoomMessagesUseCase{exec: func(_ context.Context, _ *chatApplication.ListRoomMessagesInput) (*chatApplication.ListRoomMessagesOutput, error) {
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

			router.GET("/chats/rooms/messages/:room_id", chatHTTP.NewListRoomMessagesHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
