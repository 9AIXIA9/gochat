package http_test

import (
	"context"
	"gochat/config"
	authApplication "gochat/internal/authorization/application"
	authDomain "gochat/internal/authorization/domain"
	authhttp "gochat/internal/authorization/port/http"
	ginMocks "gochat/internal/infrastructure/gin/mocks"
	"gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
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

type fakeRefreshAccessTokenUseCase struct {
	exec func(ctx context.Context, in *authApplication.RefreshAccessTokenInput) (*authApplication.RefreshAccessTokenOutput, error)
}

func (f fakeRefreshAccessTokenUseCase) Execute(ctx context.Context, in *authApplication.RefreshAccessTokenInput) (*authApplication.RefreshAccessTokenOutput, error) {
	return f.exec(ctx, in)
}

func TestRefreshAccessTokenRequest_Bind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cookie      *http.Cookie
		expectErr   bool
		expectToken authDomain.RefreshToken
	}{
		{
			name: "success",
			cookie: &http.Cookie{
				Name:  authhttp.RefreshTokenCookieKey,
				Value: "refresh-token-abc",
			},
			expectErr:   false,
			expectToken: authDomain.RefreshToken("refresh-token-abc"),
		},
		{
			name:      "missing cookie",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(responseRecorder)

			req, err := http.NewRequest(http.MethodPut, "/auth/tokens/refresh", nil)
			require.NoError(t, err)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}
			ginContext.Request = req

			request := &authhttp.RefreshAccessTokenRequest{}
			err = request.Bind(ginContext)

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectToken, request.RefreshToken)
		})
	}
}

func TestNewRefreshAccessTokenHandler(t *testing.T) {
	t.Parallel()

	fixedAccessToken := authDomain.AccessToken("access-token-xyz")
	cookieConfig := &config.Cookie{
		Path:     "/",
		Domain:   "demo.local",
		Secure:   true,
		HttpOnly: true,
	}

	tests := []struct {
		name            string
		cookie          *http.Cookie
		expectValidator func(validator *ginMocks.MockValidator)
		useCase         fakeRefreshAccessTokenUseCase
		expectResponse  *api.Response
		assertCookie    func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "success maps cookie token and refreshes cookie",
			cookie: &http.Cookie{
				Name:  authhttp.RefreshTokenCookieKey,
				Value: "refresh-token-incoming",
			},
			expectValidator: func(validator *ginMocks.MockValidator) {
				validator.EXPECT().Validate(gomock.Any(), gomock.Any()).Return("", nil)
			},
			useCase: fakeRefreshAccessTokenUseCase{exec: func(_ context.Context, in *authApplication.RefreshAccessTokenInput) (*authApplication.RefreshAccessTokenOutput, error) {
				assert.Equal(t, authDomain.RefreshToken("refresh-token-incoming"), in.RefreshToken)
				return &authApplication.RefreshAccessTokenOutput{
					AccessToken: fixedAccessToken,
					RefreshToken: authDomain.LoadRefreshToken(
						"refresh-token-refreshed",
						"user-1",
						time.Now().UTC().Add(10*time.Minute),
						0,
					),
				}, nil
			}},
			expectResponse: api.NewResponseWithData(authhttp.RefreshAccessTokenResponseData{AccessToken: fixedAccessToken}),
			assertCookie: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()
				cookies := recorder.Result().Cookies()
				require.NotEmpty(t, cookies)
				var refreshCookie *http.Cookie
				for _, cookie := range cookies {
					if cookie.Name == authhttp.RefreshTokenCookieKey {
						refreshCookie = cookie
						break
					}
				}
				require.NotNil(t, refreshCookie)
				assert.Equal(t, "refresh-token-refreshed", refreshCookie.Value)
				assert.Equal(t, cookieConfig.Path, refreshCookie.Path)
				assert.Equal(t, cookieConfig.Domain, refreshCookie.Domain)
				assert.Equal(t, cookieConfig.Secure, refreshCookie.Secure)
				assert.Equal(t, cookieConfig.HttpOnly, refreshCookie.HttpOnly)
				assert.Greater(t, refreshCookie.MaxAge, 0)
			},
		},
		{
			name:            "missing cookie returns invalid param",
			cookie:          nil,
			expectValidator: nil,
			useCase: fakeRefreshAccessTokenUseCase{exec: func(_ context.Context, _ *authApplication.RefreshAccessTokenInput) (*authApplication.RefreshAccessTokenOutput, error) {
				t.Fatalf("use case should not be called when bind fails")
				return nil, nil
			}},
			expectResponse: api.ResponseInvalidParam,
		},
		{
			name: "use case business error returns business code and message",
			cookie: &http.Cookie{
				Name:  authhttp.RefreshTokenCookieKey,
				Value: "refresh-token-incoming",
			},
			expectValidator: func(validator *ginMocks.MockValidator) {
				validator.EXPECT().Validate(gomock.Any(), gomock.Any()).Return("", nil)
			},
			useCase: fakeRefreshAccessTokenUseCase{exec: func(_ context.Context, _ *authApplication.RefreshAccessTokenInput) (*authApplication.RefreshAccessTokenOutput, error) {
				return nil, myErrors.NewBusiness("invalid refresh token")
			}},
			expectResponse: api.NewResponseWithMessage(api.CodeBusinessError, "invalid refresh token"),
		},
		{
			name: "expired refresh token in output returns server error",
			cookie: &http.Cookie{
				Name:  authhttp.RefreshTokenCookieKey,
				Value: "refresh-token-incoming",
			},
			expectValidator: func(validator *ginMocks.MockValidator) {
				validator.EXPECT().Validate(gomock.Any(), gomock.Any()).Return("", nil)
			},
			useCase: fakeRefreshAccessTokenUseCase{exec: func(_ context.Context, _ *authApplication.RefreshAccessTokenInput) (*authApplication.RefreshAccessTokenOutput, error) {
				return &authApplication.RefreshAccessTokenOutput{
					AccessToken: fixedAccessToken,
					RefreshToken: authDomain.LoadRefreshToken(
						"refresh-token-expired",
						"user-1",
						time.Now().UTC().Add(-1*time.Second),
						0,
					),
				}, nil
			}},
			expectResponse: api.ResponseServerError,
			assertCookie: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()
				cookies := recorder.Result().Cookies()
				for _, cookie := range cookies {
					assert.NotEqual(t, authhttp.RefreshTokenCookieKey, cookie.Name)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			router := httptestutil.NewTestRouter(t)

			validator := ginMocks.NewMockValidator(ctrl)
			if tt.expectValidator != nil {
				tt.expectValidator(validator)
			}

			router.PUT("/auth/tokens/refresh", authhttp.NewRefreshAccessTokenHandler(tt.useCase, validator, cookieConfig))

			req, err := http.NewRequest(http.MethodPut, "/auth/tokens/refresh", nil)
			require.NoError(t, err)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)

			if tt.assertCookie != nil {
				tt.assertCookie(t, w)
			}
		})
	}
}
