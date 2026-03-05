package http_test

import (
	"context"
	ginMocks "gochat/internal/infrastructure/gin/mocks"
	profileApplication "gochat/internal/profile/application"
	profileDomain "gochat/internal/profile/domain"
	profileDTO "gochat/internal/profile/dto"
	profileHTTP "gochat/internal/profile/port/http"
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

type fakeGetUserProfileUseCase struct {
	exec func(ctx context.Context, in *profileApplication.GetUserProfileInput) (*profileApplication.GetUserProfileOutput, error)
}

func (f fakeGetUserProfileUseCase) Execute(ctx context.Context, in *profileApplication.GetUserProfileInput) (*profileApplication.GetUserProfileOutput, error) {
	return f.exec(ctx, in)
}

func TestGetMyProfileRequest_Bind(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")

	tests := []struct {
		name     string
		userID   kernel.UserID
		expected kernel.UserID
	}{
		{
			name:     "with user id in context",
			userID:   fixedUserID,
			expected: fixedUserID,
		},
		{
			name:     "without user id in context",
			userID:   "",
			expected: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)

			req, err := http.NewRequest(http.MethodGet, "/profiles/me", nil)
			require.NoError(t, err)
			req = req.WithContext(ctxutil.WithUserID(req.Context(), tt.userID))
			ginContext.Request = req

			request := &profileHTTP.GetMyProfileRequest{}
			err = request.Bind(ginContext)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, request.UserID)
		})
	}
}

func TestNewGetMyProfileHandler(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")
	fixedSignedUpAt := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	fixedProfile := profileDomain.LoadUserProfile(
		fixedUserID,
		"Jack",
		kernel.MaleGender,
		"jack@demo.com",
		"13310001000",
		"China",
		"hello",
		fixedSignedUpAt,
	)

	tests := []struct {
		name           string
		injectUserID   kernel.UserID
		expectValidate bool
		useCase        fakeGetUserProfileUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps context user id and returns profile",
			injectUserID:   fixedUserID,
			expectValidate: true,
			useCase: fakeGetUserProfileUseCase{exec: func(_ context.Context, in *profileApplication.GetUserProfileInput) (*profileApplication.GetUserProfileOutput, error) {
				assert.Equal(t, fixedUserID, in.UserID)
				return &profileApplication.GetUserProfileOutput{Profile: fixedProfile}, nil
			}},
			expectResponse: api.NewResponseWithData(&profileHTTP.GetMyProfileResponseData{Profile: profileDTO.ToUserProfileDTO(fixedProfile)}),
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

			router.GET("/profiles/me", profileHTTP.NewGetMyProfileHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodGet, "/profiles/me", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
