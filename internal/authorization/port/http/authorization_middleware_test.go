package http_test

import (
	"context"
	"errors"
	authApplication "gochat/internal/authorization/application"
	authDomain "gochat/internal/authorization/domain"
	authHTTP "gochat/internal/authorization/port/http"
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
)

type fakeParseAccessTokenUseCase struct {
	exec func(ctx context.Context, in *authApplication.ParseAccessTokenInput) (*authApplication.ParseAccessTokenOutput, error)
}

func (f fakeParseAccessTokenUseCase) Execute(ctx context.Context, in *authApplication.ParseAccessTokenInput) (*authApplication.ParseAccessTokenOutput, error) {
	return f.exec(ctx, in)
}

func TestNewAuthorizationMiddleware(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")

	tests := []struct {
		name           string
		authorization  string
		useCase        fakeParseAccessTokenUseCase
		expectResponse *api.Response
		assertContext  bool
	}{
		{
			name:          "success with bearer token",
			authorization: "Bearer token-abc",
			useCase: fakeParseAccessTokenUseCase{exec: func(_ context.Context, in *authApplication.ParseAccessTokenInput) (*authApplication.ParseAccessTokenOutput, error) {
				assert.Equal(t, authDomain.AccessToken("token-abc"), in.AccessToken)
				return &authApplication.ParseAccessTokenOutput{UserID: fixedUserID}, nil
			}},
			expectResponse: api.ResponseSuccess,
			assertContext:  true,
		},
		{
			name: "missing authorization header",
			useCase: fakeParseAccessTokenUseCase{exec: func(_ context.Context, _ *authApplication.ParseAccessTokenInput) (*authApplication.ParseAccessTokenOutput, error) {
				t.Fatalf("use case should not be called")
				return nil, nil
			}},
			expectResponse: api.ResponseInvalidToken,
		},
		{
			name:          "invalid authorization scheme",
			authorization: "Token token-abc",
			useCase: fakeParseAccessTokenUseCase{exec: func(_ context.Context, _ *authApplication.ParseAccessTokenInput) (*authApplication.ParseAccessTokenOutput, error) {
				t.Fatalf("use case should not be called")
				return nil, nil
			}},
			expectResponse: api.ResponseInvalidToken,
		},
		{
			name:          "parse token error",
			authorization: "Bearer token-abc",
			useCase: fakeParseAccessTokenUseCase{exec: func(_ context.Context, _ *authApplication.ParseAccessTokenInput) (*authApplication.ParseAccessTokenOutput, error) {
				return nil, errors.New("invalid")
			}},
			expectResponse: api.ResponseInvalidToken,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			router := httptestutil.NewTestRouter(t)
			router.Use(authHTTP.NewAuthorizationMiddleware(tt.useCase))
			router.GET("/protected", func(c *gin.Context) {
				if tt.assertContext {
					assert.Equal(t, fixedUserID, ctxutil.UserIDFrom(c.Request.Context()))
				}
				c.JSON(http.StatusOK, api.ResponseSuccess)
			})

			req, err := http.NewRequest(http.MethodGet, "/protected", nil)
			require.NoError(t, err)
			if tt.authorization != "" {
				req.Header.Set("Authorization", tt.authorization)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
