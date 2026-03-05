package http_test

import (
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

type fakeAgreeFriendRequestUseCase struct {
	exec func(ctx context.Context, in *friendshipApplication.AgreeFriendRequestInput) (*kernel.NoOutput, error)
}

func (f fakeAgreeFriendRequestUseCase) Execute(ctx context.Context, in *friendshipApplication.AgreeFriendRequestInput) (*kernel.NoOutput, error) {
	return f.exec(ctx, in)
}

func TestAgreeFriendRequestRequest_Bind(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")
	fixedRequestID := kernel.OperationID("request-1")

	tests := []struct {
		name      string
		userID    kernel.UserID
		requestID kernel.OperationID
		want      *friendshipHTTP.AgreeFriendRequestRequest
	}{
		{
			name:      "success",
			userID:    fixedUserID,
			requestID: fixedRequestID,
			want: &friendshipHTTP.AgreeFriendRequestRequest{
				UserID:    fixedUserID,
				RequestID: fixedRequestID,
			},
		},
		{
			name:      "missing request id",
			userID:    fixedUserID,
			requestID: "",
			want: &friendshipHTTP.AgreeFriendRequestRequest{
				UserID:    fixedUserID,
				RequestID: "",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)

			req, err := http.NewRequest(http.MethodPut, "/friendship-requests/request-1/agree", nil)
			require.NoError(t, err)
			req = req.WithContext(ctxutil.WithUserID(req.Context(), tt.userID))
			ginContext.Request = req

			if tt.requestID != "" {
				ginContext.Params = gin.Params{{Key: "request_id", Value: string(tt.requestID)}}
			}

			request := &friendshipHTTP.AgreeFriendRequestRequest{}
			err = request.Bind(ginContext)
			require.NoError(t, err)
			assert.Equal(t, tt.want, request)
		})
	}
}

func TestNewAgreeFriendRequestHandler(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")
	fixedRequestID := kernel.OperationID("request-1")

	tests := []struct {
		name           string
		injectUserID   kernel.UserID
		path           string
		expectValidate bool
		useCase        fakeAgreeFriendRequestUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps request to input",
			injectUserID:   fixedUserID,
			path:           "/friendship-requests/request-1/agree",
			expectValidate: true,
			useCase: fakeAgreeFriendRequestUseCase{exec: func(_ context.Context, in *friendshipApplication.AgreeFriendRequestInput) (*kernel.NoOutput, error) {
				assert.Equal(t, fixedUserID, in.UserID)
				assert.Equal(t, fixedRequestID, in.RequestID)
				return nil, nil
			}},
			expectResponse: api.ResponseSuccess,
		},
		{
			name:           "missing user id returns invalid param",
			injectUserID:   "",
			path:           "/friendship-requests/request-1/agree",
			expectValidate: true,
			useCase: fakeAgreeFriendRequestUseCase{exec: func(_ context.Context, _ *friendshipApplication.AgreeFriendRequestInput) (*kernel.NoOutput, error) {
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

			router.PUT("/friendship-requests/:request_id/agree", friendshipHTTP.NewAgreeFriendRequestHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodPut, tt.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
