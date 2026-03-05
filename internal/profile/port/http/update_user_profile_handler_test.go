package http_test

import (
	"bytes"
	"context"
	ginMocks "gochat/internal/infrastructure/gin/mocks"
	profileApplication "gochat/internal/profile/application"
	profileHTTP "gochat/internal/profile/port/http"
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

type fakeUpdateUserProfileUseCase struct {
	exec func(ctx context.Context, in *profileApplication.UpdateUserProfileInput) (*kernel.NoOutput, error)
}

func (f fakeUpdateUserProfileUseCase) Execute(ctx context.Context, in *profileApplication.UpdateUserProfileInput) (*kernel.NoOutput, error) {
	return f.exec(ctx, in)
}

func TestUpdateUserProfileRequest_Bind(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")

	tests := []struct {
		name      string
		userID    kernel.UserID
		body      string
		expectErr bool
		want      *profileHTTP.UpdateUserProfileRequest
	}{
		{
			name:      "success",
			userID:    fixedUserID,
			body:      `{"name":"Jack","gender":1,"email":"jack@demo.com","phone_number":"13310001000","address":"China","sign":"hello"}`,
			expectErr: false,
			want: &profileHTTP.UpdateUserProfileRequest{
				UserID:      fixedUserID,
				Name:        "Jack",
				Gender:      kernel.MaleGender,
				Email:       kernel.Email("jack@demo.com"),
				PhoneNumber: kernel.PhoneNumber("13310001000"),
				Address:     kernel.Address("China"),
				Sign:        "hello",
			},
		},
		{
			name:      "malformed json",
			userID:    fixedUserID,
			body:      `{"name":`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)

			req, err := http.NewRequest(http.MethodPut, "/profiles/me", bytes.NewBufferString(tt.body))
			require.NoError(t, err)
			req = req.WithContext(ctxutil.WithUserID(req.Context(), tt.userID))
			req.Header.Set("Content-Type", "application/json")
			ginContext.Request = req

			request := &profileHTTP.UpdateUserProfileRequest{}
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

func TestNewUpdateUserProfileHandler(t *testing.T) {
	t.Parallel()

	fixedUserID := kernel.UserID("user-1")
	body := `{"name":"Jack","gender":1,"email":"jack@demo.com","phone_number":"13310001000","address":"China","sign":"hello"}`

	tests := []struct {
		name           string
		injectUserID   kernel.UserID
		body           string
		expectValidate bool
		useCase        fakeUpdateUserProfileUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps request to input",
			injectUserID:   fixedUserID,
			body:           body,
			expectValidate: true,
			useCase: fakeUpdateUserProfileUseCase{exec: func(_ context.Context, in *profileApplication.UpdateUserProfileInput) (*kernel.NoOutput, error) {
				assert.Equal(t, fixedUserID, in.UserID)
				assert.Equal(t, "Jack", in.Name)
				assert.Equal(t, kernel.MaleGender, in.Gender)
				assert.Equal(t, kernel.Email("jack@demo.com"), in.Email)
				assert.Equal(t, kernel.PhoneNumber("13310001000"), in.PhoneNumber)
				assert.Equal(t, kernel.Address("China"), in.Address)
				assert.Equal(t, "hello", in.Sign)
				return nil, nil
			}},
			expectResponse: api.ResponseSuccess,
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

			router.PUT("/profiles/me", profileHTTP.NewUpdateUserProfileHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodPut, "/profiles/me", bytes.NewBufferString(tt.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
