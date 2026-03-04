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

func TestGetUserProfileRequest_Bind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		params    gin.Params
		expectErr bool
		wantID    kernel.UserID
	}{
		{
			name:      "success",
			params:    gin.Params{{Key: "user_id", Value: "user-1"}},
			expectErr: false,
			wantID:    kernel.UserID("user-1"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Params = tt.params
			req, err := http.NewRequest(http.MethodGet, "/profiles/users/user-1", nil)
			require.NoError(t, err)
			ginContext.Request = req

			request := &profileHTTP.GetUserProfileRequest{}
			err = request.Bind(ginContext)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantID, request.UserID)
		})
	}
}

func TestNewGetUserProfileHandler(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")
	fixedProfile := profileDomain.LoadUserProfile(
		fixedUserID,
		"Jack",
		kernel.MaleGender,
		"jack@demo.com",
		"13310001000",
		"China",
		"hello",
		time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC),
	)

	tests := []struct {
		name           string
		path           string
		expectValidate bool
		useCase        fakeGetUserProfileUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps uri user id and returns profile",
			path:           "/profiles/users/user-1",
			expectValidate: true,
			useCase: fakeGetUserProfileUseCase{exec: func(_ context.Context, in *profileApplication.GetUserProfileInput) (*profileApplication.GetUserProfileOutput, error) {
				assert.Equal(t, fixedUserID, in.UserID)
				return &profileApplication.GetUserProfileOutput{Profile: fixedProfile}, nil
			}},
			expectResponse: api.NewResponseWithData(&profileHTTP.GetUserProfileResponseData{Profile: profileDTO.ToUserProfileDTO(fixedProfile)}),
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

			router.GET("/profiles/users/:user_id", profileHTTP.NewGetUserProfileHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
