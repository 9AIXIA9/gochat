package http_test

import (
	"bytes"
	"context"
	authApplication "gochat/internal/authorization/application"
	authDomain "gochat/internal/authorization/domain"
	authhttp "gochat/internal/authorization/port/http"
	ginMocks "gochat/internal/infrastructure/gin/mocks"
	"gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	httptestutil "gochat/pkg/httptest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type fakeSignUpUseCase struct {
	exec func(ctx context.Context, in *authApplication.SignUpInput) (*authApplication.SignUpOutput, error)
}

func (f fakeSignUpUseCase) Execute(ctx context.Context, in *authApplication.SignUpInput) (*authApplication.SignUpOutput, error) {
	return f.exec(ctx, in)
}

func TestSignUpRequest_Bind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		body      string
		expectErr bool
		want      *authhttp.SignUpRequest
	}{
		{
			name:      "success",
			body:      `{"email":"youremail@demo.com","password":"your-password"}`,
			expectErr: false,
			want: &authhttp.SignUpRequest{
				Email:    kernel.Email("youremail@demo.com"),
				Password: authDomain.Password("your-password"),
			},
		},
		{
			name:      "malformed json",
			body:      `{"email":`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(responseRecorder)

			req, err := http.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(tt.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			ginContext.Request = req

			request := &authhttp.SignUpRequest{}
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

func TestNewSignUpHandler(t *testing.T) {
	t.Parallel()

	fixedEmail := kernel.Email("youremail@demo.com")
	fixedPassword := authDomain.Password("your-password")
	fixedUserNumber := kernel.UserNumber("2004426295315795968")

	tests := []struct {
		name            string
		body            string
		expectValidator func(validator *ginMocks.MockValidator)
		useCase         fakeSignUpUseCase
		expectResponse  *api.Response
	}{
		{
			name: "success maps request to input and returns user_number",
			body: `{"email":"youremail@demo.com","password":"your-password"}`,
			expectValidator: func(validator *ginMocks.MockValidator) {
				validator.EXPECT().Validate(gomock.Any(), gomock.Any()).Return("", nil)
			},
			useCase: fakeSignUpUseCase{exec: func(_ context.Context, in *authApplication.SignUpInput) (*authApplication.SignUpOutput, error) {
				assert.Equal(t, fixedEmail, in.Email)
				assert.Equal(t, fixedPassword, in.Password)
				return &authApplication.SignUpOutput{UserNumber: fixedUserNumber}, nil
			}},
			expectResponse: api.NewResponseWithData(&authhttp.SignUpResponseData{UserNumber: fixedUserNumber}),
		},
		{
			name:            "invalid json returns invalid param",
			body:            `{"email":`,
			expectValidator: nil,
			useCase: fakeSignUpUseCase{exec: func(_ context.Context, _ *authApplication.SignUpInput) (*authApplication.SignUpOutput, error) {
				t.Fatalf("use case should not be called when bind fails")
				return nil, nil
			}},
			expectResponse: api.ResponseInvalidParam,
		},
		{
			name: "use case business error returns business code and message",
			body: `{"email":"youremail@demo.com","password":"your-password"}`,
			expectValidator: func(validator *ginMocks.MockValidator) {
				validator.EXPECT().Validate(gomock.Any(), gomock.Any()).Return("", nil)
			},
			useCase: fakeSignUpUseCase{exec: func(_ context.Context, _ *authApplication.SignUpInput) (*authApplication.SignUpOutput, error) {
				return nil, myErrors.NewBusiness("email is already used")
			}},
			expectResponse: api.NewResponseWithMessage(api.CodeBusinessError, "email is already used"),
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

			router.POST("/auth/sign-up", authhttp.NewSignUpHandler(tt.useCase, validator))

			httptestutil.ExecuteRouteTest(t, router, &httptestutil.TestCase{
				Path:     "/auth/sign-up",
				Method:   http.MethodPost,
				Body:     bytes.NewBufferString(tt.body),
				Response: tt.expectResponse,
			})
		})
	}
}
