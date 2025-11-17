package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/validator"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/application/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LeaveRoomRequest struct {
	UserID     kernel.UserID     `json:"-" validate:"required"`
	RoomNumber kernel.RoomNumber `json:"room_number" validate:"required"`
}

func (r *LeaveRoomRequest) Bind(ginContext *gin.Context) error {
	id := ginutils.GetUserID(ginContext)
	r.UserID = id
	return ginContext.ShouldBind(r)
}

func NewLeaveRoomHandler(useCase usecase.LeaveRoomUseCase, validator *validator.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *LeaveRoomRequest) *usecase.LeaveRoomInput {
			return &usecase.LeaveRoomInput{
				UserID:     request.UserID,
				RoomNumber: request.RoomNumber,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.Response(ginContext, sharedHttp.ResponseSuccess)
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "input is empty"))
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "room not found"))
			case errors.Is(err, myErrors.ErrOwnerCantLeave):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeOwnerCantLeave, "owner can't leave the room"))

			default:
				zap.L().Error("leave room handler failed", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
