package http_test

import (
	"context"
	ginMocks "gochat/internal/infrastructure/gin/mocks"
	roomshipApplication "gochat/internal/roomship/application"
	roomshipDomain "gochat/internal/roomship/domain"
	roomshipDTO "gochat/internal/roomship/dto"
	roomshipHTTP "gochat/internal/roomship/port/http"
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

type fakeListRoomMembersUseCase struct {
	exec func(ctx context.Context, in *roomshipApplication.ListRoomMembersInput) (*roomshipApplication.ListRoomMembersOutput, error)
}

func (f fakeListRoomMembersUseCase) Execute(ctx context.Context, in *roomshipApplication.ListRoomMembersInput) (*roomshipApplication.ListRoomMembersOutput, error) {
	return f.exec(ctx, in)
}

func TestListRoomMembersRequest_Bind(t *testing.T) {
	t.Parallel()

	fixedRoomID := kernel.RoomID("room-1")

	tests := []struct {
		name   string
		roomID kernel.RoomID
		want   *roomshipHTTP.ListRoomMembersRequest
	}{
		{
			name:   "success",
			roomID: fixedRoomID,
			want:   &roomshipHTTP.ListRoomMembersRequest{RoomID: fixedRoomID},
		},
		{
			name:   "missing room id",
			roomID: "",
			want:   &roomshipHTTP.ListRoomMembersRequest{RoomID: ""},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)

			req, err := http.NewRequest(http.MethodGet, "/rooms/room-1/members", nil)
			require.NoError(t, err)
			ginContext.Request = req

			if tt.roomID != "" {
				ginContext.Params = gin.Params{{Key: "room_id", Value: string(tt.roomID)}}
			}

			request := &roomshipHTTP.ListRoomMembersRequest{}
			err = request.Bind(ginContext)
			require.NoError(t, err)
			assert.Equal(t, tt.want, request)
		})
	}
}

func TestNewListRoomMembersHandler(t *testing.T) {
	t.Parallel()

	fixedRoomID := kernel.RoomID("room-1")
	fixedCreatedAt := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	fixedRoomship := roomshipDomain.LoadRoomship(
		"roomship-1",
		fixedRoomID,
		"user-1",
		roomshipDomain.MemberRole,
		fixedCreatedAt,
	)

	tests := []struct {
		name           string
		path           string
		expectValidate bool
		useCase        fakeListRoomMembersUseCase
		expectResponse *api.Response
	}{
		{
			name:           "success maps request to input and returns members",
			path:           "/rooms/room-1/members",
			expectValidate: true,
			useCase: fakeListRoomMembersUseCase{exec: func(_ context.Context, in *roomshipApplication.ListRoomMembersInput) (*roomshipApplication.ListRoomMembersOutput, error) {
				assert.Equal(t, fixedRoomID, in.RoomID)
				return &roomshipApplication.ListRoomMembersOutput{Roomships: []*roomshipDomain.Roomship{fixedRoomship}}, nil
			}},
			expectResponse: api.NewResponseWithData(&roomshipHTTP.ListRoomMembersResponseData{Roomships: roomshipDTO.ToRoomshipDTOs([]*roomshipDomain.Roomship{fixedRoomship})}),
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

			router.GET("/rooms/:room_id/members", roomshipHTTP.NewListRoomMembersHandler(tt.useCase, validator))

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectResponse.String(), w.Body.String())
			assert.Equal(t, tt.expectResponse.Code.ToHTTPCode(), w.Code)
		})
	}
}
