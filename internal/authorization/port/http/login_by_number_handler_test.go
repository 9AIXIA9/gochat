package http_test

import (
	"bytes"
	"context"
	"gochat/config"
	authApplication "gochat/internal/authorization/application"
	authDomain "gochat/internal/authorization/domain"
	authhttp "gochat/internal/authorization/port/http"
	ginMocks "gochat/internal/infrastructure/gin/mocks"
	"gochat/internal/shared/api"
	"gochat/internal/shared/kernel"
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

type fakeLoginByNumberUseCase struct {
	exec func(ctx context.Context, in *authApplication.LoginByNumberInput) (*authApplication.LoginByNumberOutput, error)
}

func (f fakeLoginByNumberUseCase) Execute(ctx context.Context, in *authApplication.LoginByNumberInput) (*authApplication.LoginByNumberOutput, error) {
	return f.exec(ctx, in)
}

func TestLoginByNumberRequest_Bind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		body      string
		expectErr bool
		want      *authhttp.LoginByNumberRequest
	}{
		{
			name:      "success",
			body:      `{"number":"2004426295315795968","password":"your-password"}`,
			expectErr: false,
			want: &authhttp.LoginByNumberRequest{
				Number:   kernel.UserNumber("2004426295315795968"),
				Password: authDomain.Password("your-password"),
			},
		},
		{
			name:      "malformed json",
			body:      `{"number":`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(responseRecorder)

			req, err := http.NewRequest(http.MethodPost, "/auth/login/user_number", bytes.NewBufferString(tt.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			ginContext.Request = req

			request := &authhttp.LoginByNumberRequest{}
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

func TestNewLoginByNumberHandler(t *testing.T) {
	t.Parallel()

	fixedNumber := kernel.UserNumber("2004426295315795968")
	fixedPassword := authDomain.Password("your-password")
	fixedAccessToken := authDomain.AccessToken("access-token-xyz")
	cookieConfig := &config.Cookie{
		Path:     "/",
		Domain:   "demo.local",
		Secure:   true,
		HttpOnly: true,
	}

	tests := []struct {
		name           string
		body           string
		useCase        fakeLoginByNumberUseCase
		expectValidate bool
		expectResponse *api.Response
		assertCookie   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "success maps request to input and sets refresh cookie",
			body: `{"number":"2004426295315795968","password":"your-password"}`,
			useCase: fakeLoginByNumberUseCase{exec: func(_ context.Context, in *authApplication.LoginByNumberInput) (*authApplication.LoginByNumberOutput, error) {
				assert.Equal(t, fixedNumber, in.Number)
				assert.Equal(t, fixedPassword, in.Password)
				return &authApplication.LoginByNumberOutput{
					AccessToken: fixedAccessToken,
					RefreshToken: authDomain.LoadRefreshToken(
						"refresh-token-abc",
						"user-1",
						time.Now().UTC().Add(10*time.Minute),
						0,
					),
				}, nil
			}},
			expectValidate: true,
			expectResponse: api.NewResponseWithData(&authhttp.LoginByNumberResponseData{AccessToken: fixedAccessToken}),
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
				assert.Equal(t, "refresh-token-abc", refreshCookie.Value)
				assert.Equal(t, cookieConfig.Path, refreshCookie.Path)
				assert.Equal(t, cookieConfig.Domain, refreshCookie.Domain)
				assert.Equal(t, cookieConfig.Secure, refreshCookie.Secure)
				assert.Equal(t, cookieConfig.HttpOnly, refreshCookie.HttpOnly)
				assert.Greater(t, refreshCookie.MaxAge, 0)
			},
		},
		{
			name: "expired refresh token in output returns server error",
			body: `{"number":"2004426295315795968","password":"your-password"}`,
			useCase: fakeLoginByNumberUseCase{exec: func(_ context.Context, _ *authApplication.LoginByNumberInput) (*authApplication.LoginByNumberOutput, error) {
				return &authApplication.LoginByNumberOutput{
					AccessToken: fixedAccessToken,
					RefreshToken: authDomain.LoadRefreshToken(
						"refresh-token-expired",
						"user-1",
						time.Now().UTC().Add(-1*time.Second),
						0,
					),
				}, nil
			}},
			expectValidate: true,
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

			if tt.expectValidate {
				validator.EXPECT().Validate(gomock.Any(), gomock.Any()).Return("", nil)
			}

			router.POST("/auth/login/user_number", authhttp.NewLoginByNumberHandler(tt.useCase, validator, cookieConfig))

			req, err := http.NewRequest(http.MethodPost, "/auth/login/user_number", bytes.NewBufferString(tt.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

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
