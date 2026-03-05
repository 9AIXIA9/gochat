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

type fakeGetRoomProfileUseCase struct {
	exec func(ctx context.Context, in *profileApplication.GetRoomProfileInput) (*profileApplication.GetRoomProfileOutput, error)
}

func (f fakeGetRoomProfileUseCase) Execute(ctx context.Context, in *profileApplication.GetRoomProfileInput) (*profileApplication.GetRoomProfileOutput, error) {
	return f.exec(ctx, in)
}

func TestGetRoomProfileRequest_Bind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		params    gin.Params
		expectErr bool
		wantID    kernel.RoomID
	}{
		{
			name:      "success",
			params:    gin.Params{{Key: "room_id", Value: "room-1"}},
			expectErr: false,
			wantID:    kernel.RoomID("room-1"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Params = tt.params
			req, err := http.NewRequest(http.MethodGet, "/profiles/rooms/room-1", nil)
			require.NoError(t, err)
			ginContext.Request = req

			request := &profileHTTP.GetRoomProfileRequest{}
			err = request.Bind(ginContext)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantID, request.RoomID)
		})
	}
}

func TestNewGetRoomProfileHandler(t *testing.T) {
	t.Parallel()

	fixedRoomID := kernel.RoomID("room-1")
	fixedProfile := profileDomain.LoadRoomProfile(
		fixedRoomID,
		"family",
		"our room",
		time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC),
	)

	tests := []struct {
		name           string
		path           string
		expectValidate bool
		useCase        fakeGetRoomProfileUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps uri room id and returns profile",
			path:           "/profiles/rooms/room-1",
			expectValidate: true,
			useCase: fakeGetRoomProfileUseCase{exec: func(_ context.Context, in *profileApplication.GetRoomProfileInput) (*profileApplication.GetRoomProfileOutput, error) {
				assert.Equal(t, fixedRoomID, in.RoomID)
				return &profileApplication.GetRoomProfileOutput{Profile: fixedProfile}, nil
			}},
			expectResponse: api.NewResponseWithData(&profileHTTP.GetRoomProfileResponseData{Profile: profileDTO.ToRoomProfileDTO(fixedProfile)}),
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

			router.GET("/profiles/rooms/:room_id", profileHTTP.NewGetRoomProfileHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
