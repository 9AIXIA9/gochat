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

type fakeListPrivateMessagesUseCase struct {
	exec func(ctx context.Context, in *chatApplication.ListPrivateMessagesInput) (*chatApplication.ListPrivateMessagesOutput, error)
}

func (f fakeListPrivateMessagesUseCase) Execute(ctx context.Context, in *chatApplication.ListPrivateMessagesInput) (*chatApplication.ListPrivateMessagesOutput, error) {
	return f.exec(ctx, in)
}

func TestListPrivateMessagesRequest_Bind(t *testing.T) {
	t.Parallel()

	fixedOperatorID := kernel.UserID("user-1")

	tests := []struct {
		name      string
		userID    kernel.UserID
		path      string
		params    gin.Params
		expectErr bool
		want      *chatHTTP.ListPrivateMessagesRequest
	}{
		{
			name:      "success with query and uri",
			userID:    fixedOperatorID,
			path:      "/chats/private-messages/user-2?base_id=message-9&limit=30",
			params:    gin.Params{{Key: "user_id", Value: "user-2"}},
			expectErr: false,
			want: &chatHTTP.ListPrivateMessagesRequest{
				OperatorID: fixedOperatorID,
				UserID:     kernel.UserID("user-2"),
				BaseID:     kernel.MessageID("message-9"),
				Limit:      30,
			},
		},
		{
			name:      "default limit",
			userID:    fixedOperatorID,
			path:      "/chats/private-messages/user-2",
			params:    gin.Params{{Key: "user_id", Value: "user-2"}},
			expectErr: false,
			want: &chatHTTP.ListPrivateMessagesRequest{
				OperatorID: fixedOperatorID,
				UserID:     kernel.UserID("user-2"),
				BaseID:     "",
				Limit:      20,
			},
		},
		{
			name:      "invalid limit",
			userID:    fixedOperatorID,
			path:      "/chats/private-messages/user-2?limit=bad",
			params:    gin.Params{{Key: "user_id", Value: "user-2"}},
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

			request := &chatHTTP.ListPrivateMessagesRequest{}
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

func TestNewListPrivateMessagesHandler(t *testing.T) {
	t.Parallel()

	fixedOperatorID := kernel.UserID("user-1")
	fixedUserID := kernel.UserID("user-2")
	fixedBaseID := kernel.MessageID("message-9")
	fixedSentAt := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	fixedMessage := chatDomain.LoadPrivateMessage(
		"message-10",
		fixedOperatorID,
		fixedUserID,
		"hello",
		chatDomain.MessageStateDelivered,
		fixedSentAt,
	)

	tests := []struct {
		name           string
		injectUserID   kernel.UserID
		path           string
		expectValidate bool
		useCase        fakeListPrivateMessagesUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps request to input and returns messages",
			injectUserID:   fixedOperatorID,
			path:           "/chats/private-messages/user-2?base_id=message-9&limit=30",
			expectValidate: true,
			useCase: fakeListPrivateMessagesUseCase{exec: func(_ context.Context, in *chatApplication.ListPrivateMessagesInput) (*chatApplication.ListPrivateMessagesOutput, error) {
				assert.Equal(t, fixedOperatorID, in.OperatorID)
				assert.Equal(t, fixedUserID, in.UserID)
				assert.Equal(t, fixedBaseID, in.BaseID)
				assert.Equal(t, 30, in.Limit)
				return &chatApplication.ListPrivateMessagesOutput{PrivateMessages: []*chatDomain.PrivateMessage{fixedMessage}}, nil
			}},
			expectResponse: api.NewResponseWithData(&chatHTTP.ListPrivateMessagesResponseData{PrivateMessages: chatDTO.ToPrivateMessageDTOs([]*chatDomain.PrivateMessage{fixedMessage})}),
		},
		{
			name:           "invalid query returns invalid param",
			injectUserID:   fixedOperatorID,
			path:           "/chats/private-messages/user-2?limit=bad",
			expectValidate: false,
			useCase: fakeListPrivateMessagesUseCase{exec: func(_ context.Context, _ *chatApplication.ListPrivateMessagesInput) (*chatApplication.ListPrivateMessagesOutput, error) {
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

			router.GET("/chats/private-messages/:user_id", chatHTTP.NewListPrivateMessagesHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
