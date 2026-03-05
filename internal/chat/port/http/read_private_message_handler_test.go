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

type fakeReadPrivateMessagesUseCase struct {
	exec func(ctx context.Context, in *chatApplication.ReadPrivateMessagesInput) (*kernel.NoOutput, error)
}

func (f fakeReadPrivateMessagesUseCase) Execute(ctx context.Context, in *chatApplication.ReadPrivateMessagesInput) (*kernel.NoOutput, error) {
	return f.exec(ctx, in)
}

func TestReadPrivateMessagesRequest_Bind(t *testing.T) {
	t.Parallel()

	fixedRecipientID := kernel.UserID("user-1")

	tests := []struct {
		name      string
		userID    kernel.UserID
		body      string
		expectErr bool
		want      *chatHTTP.ReadPrivateMessagesRequest
	}{
		{
			name:      "success",
			userID:    fixedRecipientID,
			body:      `{"sender_id":"user-2"}`,
			expectErr: false,
			want: &chatHTTP.ReadPrivateMessagesRequest{
				SenderID:    kernel.UserID("user-2"),
				RecipientID: fixedRecipientID,
			},
		},
		{
			name:      "malformed json",
			userID:    fixedRecipientID,
			body:      `{"sender_id":`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)

			req, err := http.NewRequest(http.MethodPut, "/chats/private-messages/read", bytes.NewBufferString(tt.body))
			require.NoError(t, err)
			req = req.WithContext(ctxutil.WithUserID(req.Context(), tt.userID))
			req.Header.Set("Content-Type", "application/json")
			ginContext.Request = req

			request := &chatHTTP.ReadPrivateMessagesRequest{}
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

func TestNewReadPrivateMessagesHandler(t *testing.T) {
	t.Parallel()

	fixedRecipientID := kernel.UserID("user-1")
	body := `{"sender_id":"user-2"}`

	tests := []struct {
		name           string
		injectUserID   kernel.UserID
		body           string
		expectValidate bool
		useCase        fakeReadPrivateMessagesUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps request to input",
			injectUserID:   fixedRecipientID,
			body:           body,
			expectValidate: true,
			useCase: fakeReadPrivateMessagesUseCase{exec: func(_ context.Context, in *chatApplication.ReadPrivateMessagesInput) (*kernel.NoOutput, error) {
				assert.Equal(t, kernel.UserID("user-2"), in.SenderID)
				assert.Equal(t, fixedRecipientID, in.RecipientID)
				return nil, nil
			}},
			expectResponse: api.ResponseSuccess,
		},
		{
			name:           "malformed json returns invalid param",
			injectUserID:   fixedRecipientID,
			body:           `{"sender_id":`,
			expectValidate: false,
			useCase: fakeReadPrivateMessagesUseCase{exec: func(_ context.Context, _ *chatApplication.ReadPrivateMessagesInput) (*kernel.NoOutput, error) {
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

			router.PUT("/chats/private-messages/read", chatHTTP.NewReadPrivateMessagesHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodPut, "/chats/private-messages/read", bytes.NewBufferString(tt.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
