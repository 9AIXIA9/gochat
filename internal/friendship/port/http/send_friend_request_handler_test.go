package http_test

import (
	"bytes"
	"context"
	friendshipApplication "gochat/internal/friendship/application"
	friendshipHTTP "gochat/internal/friendship/port/http"
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

type fakeSendFriendRequestUseCase struct {
	exec func(ctx context.Context, in *friendshipApplication.SendFriendRequestInput) (*kernel.NoOutput, error)
}

func (f fakeSendFriendRequestUseCase) Execute(ctx context.Context, in *friendshipApplication.SendFriendRequestInput) (*kernel.NoOutput, error) {
	return f.exec(ctx, in)
}

func TestSendFriendRequestRequest_Bind(t *testing.T) {
	t.Parallel()

	fixedFromID := kernel.UserID("user-1")

	tests := []struct {
		name      string
		userID    kernel.UserID
		body      string
		expectErr bool
		want      *friendshipHTTP.SendFriendRequestRequest
	}{
		{
			name:      "success",
			userID:    fixedFromID,
			body:      `{"to_id":"user-2","content":"hello"}`,
			expectErr: false,
			want: &friendshipHTTP.SendFriendRequestRequest{
				FromID:  fixedFromID,
				ToID:    kernel.UserID("user-2"),
				Content: "hello",
			},
		},
		{
			name:      "malformed json",
			userID:    fixedFromID,
			body:      `{"to_id":`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)

			req, err := http.NewRequest(http.MethodPost, "/friendship-requests", bytes.NewBufferString(tt.body))
			require.NoError(t, err)
			req = req.WithContext(ctxutil.WithUserID(req.Context(), tt.userID))
			req.Header.Set("Content-Type", "application/json")
			ginContext.Request = req

			request := &friendshipHTTP.SendFriendRequestRequest{}
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

func TestNewSendFriendRequestHandler(t *testing.T) {
	t.Parallel()

	fixedFromID := kernel.UserID("user-1")
	body := `{"to_id":"user-2","content":"hello"}`

	tests := []struct {
		name           string
		injectUserID   kernel.UserID
		body           string
		expectValidate bool
		useCase        fakeSendFriendRequestUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps request to input",
			injectUserID:   fixedFromID,
			body:           body,
			expectValidate: true,
			useCase: fakeSendFriendRequestUseCase{exec: func(_ context.Context, in *friendshipApplication.SendFriendRequestInput) (*kernel.NoOutput, error) {
				assert.Equal(t, fixedFromID, in.FromID)
				assert.Equal(t, kernel.UserID("user-2"), in.ToID)
				assert.Equal(t, "hello", in.Content)
				return nil, nil
			}},
			expectResponse: api.ResponseSuccess,
		},
		{
			name:           "malformed json returns invalid param",
			injectUserID:   fixedFromID,
			body:           `{"to_id":`,
			expectValidate: false,
			useCase: fakeSendFriendRequestUseCase{exec: func(_ context.Context, _ *friendshipApplication.SendFriendRequestInput) (*kernel.NoOutput, error) {
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

			router.POST("/friendship-requests", friendshipHTTP.NewSendFriendRequestHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodPost, "/friendship-requests", bytes.NewBufferString(tt.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
