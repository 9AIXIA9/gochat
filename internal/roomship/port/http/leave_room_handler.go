package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LeaveRoomRequest struct {
	UserID kernel.UserID `json:"-" validate:"required"`
	RoomID kernel.RoomID `uri:"room_id" validate:"required"`
}

func (r *LeaveRoomRequest) Bind(ginContext *gin.Context) error {
	userID := ginutils.GetUserID(ginContext)
	r.UserID = userID
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

func NewLeaveRoomHandler(useCase application.LeaveRoomUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *LeaveRoomRequest) *application.LeaveRoomInput {
			return &application.LeaveRoomInput{
				UserID: request.UserID,
				RoomID: request.RoomID,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.ResponseSuccess(ginContext)
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "input is empty")
			case errors.Is(err, domain.ErrOwnerCantLeave):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "owner can't leave the room")
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.ResponseSuccess(ginContext)
			default:
				zap.L().Error("LeaveRoomHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.CodeServerError)
			}
		},
		5*time.Second,
	)
}
