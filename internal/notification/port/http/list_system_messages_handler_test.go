package http_test

import (
	"context"
	ginMocks "gochat/internal/infrastructure/gin/mocks"
	notificationApplication "gochat/internal/notification/application"
	notificationDomain "gochat/internal/notification/domain"
	notificationDTO "gochat/internal/notification/dto"
	notificationHTTP "gochat/internal/notification/port/http"
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

type fakeListSystemMessagesUseCase struct {
	exec func(ctx context.Context, in *notificationApplication.ListSystemMessagesInput) (*notificationApplication.ListSystemMessagesOutput, error)
}

func (f fakeListSystemMessagesUseCase) Execute(ctx context.Context, in *notificationApplication.ListSystemMessagesInput) (*notificationApplication.ListSystemMessagesOutput, error) {
	return f.exec(ctx, in)
}

func TestListSystemMessagesRequest_Bind(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")

	tests := []struct {
		name      string
		userID    kernel.UserID
		path      string
		expectErr bool
		want      *notificationHTTP.ListSystemMessagesRequest
	}{
		{
			name:      "success with query",
			userID:    fixedUserID,
			path:      "/notifications/system-messages?base_id=message-9&limit=30",
			expectErr: false,
			want: &notificationHTTP.ListSystemMessagesRequest{
				UserID: fixedUserID,
				BaseID: kernel.MessageID("message-9"),
				Limit:  30,
			},
		},
		{
			name:      "default limit",
			userID:    fixedUserID,
			path:      "/notifications/system-messages",
			expectErr: false,
			want: &notificationHTTP.ListSystemMessagesRequest{
				UserID: fixedUserID,
				BaseID: "",
				Limit:  20,
			},
		},
		{
			name:      "invalid limit",
			userID:    fixedUserID,
			path:      "/notifications/system-messages?limit=bad",
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

			request := &notificationHTTP.ListSystemMessagesRequest{}
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

func TestNewListSystemMessagesHandler(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")
	fixedBaseID := kernel.MessageID("message-9")
	fixedSentAt := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	fixedMessage := notificationDomain.LoadSystemMessage(
		"message-10",
		fixedUserID,
		notificationDomain.MessageStateDelivered,
		"system notice",
		fixedSentAt,
	)

	tests := []struct {
		name           string
		injectUserID   kernel.UserID
		path           string
		expectValidate bool
		useCase        fakeListSystemMessagesUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps request to input and returns system messages",
			injectUserID:   fixedUserID,
			path:           "/notifications/system-messages?base_id=message-9&limit=30",
			expectValidate: true,
			useCase: fakeListSystemMessagesUseCase{exec: func(_ context.Context, in *notificationApplication.ListSystemMessagesInput) (*notificationApplication.ListSystemMessagesOutput, error) {
				assert.Equal(t, fixedUserID, in.UserID)
				assert.Equal(t, fixedBaseID, in.BaseID)
				assert.Equal(t, 30, in.Limit)
				return &notificationApplication.ListSystemMessagesOutput{SystemMessages: []*notificationDomain.SystemMessage{fixedMessage}}, nil
			}},
			expectResponse: api.NewResponseWithData(&notificationHTTP.ListSystemMessagesResponseData{SystemMessages: notificationDTO.ToSystemMessageDTOs([]*notificationDomain.SystemMessage{fixedMessage})}),
		},
		{
			name:           "invalid query returns invalid param",
			injectUserID:   fixedUserID,
			path:           "/notifications/system-messages?limit=bad",
			expectValidate: false,
			useCase: fakeListSystemMessagesUseCase{exec: func(_ context.Context, _ *notificationApplication.ListSystemMessagesInput) (*notificationApplication.ListSystemMessagesOutput, error) {
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

			router.GET("/notifications/system-messages", notificationHTTP.NewListSystemMessagesHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
