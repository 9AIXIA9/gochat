package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/profile/application"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UpdateRoomProfileRequest struct {
	UserID       kernel.UserID `json:"-" validate:"required"`
	RoomID       kernel.RoomID `json:"room_id" validate:"required"`
	Name         string        `json:"name"`
	Introduction string        `json:"introduction"`
}

func (r *UpdateRoomProfileRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	// 绑定查询参数
	if err := ginContext.ShouldBind(r); err != nil {
		return err
	}
	return nil
}

func NewUpdateRoomProfileHandler(useCase application.UpdateRoomProfileUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *UpdateRoomProfileRequest) *application.UpdateRoomProfileInput {
			return &application.UpdateRoomProfileInput{
				UserID:       request.UserID,
				RoomID:       request.RoomID,
				Name:         request.Name,
				Introduction: request.Introduction,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.ResponseSuccess(ginContext)
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "input is empty")
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeNotFound, "room profile not found")
			default:
				zap.L().Error("UpdateRoomProfileHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.CodeServerError)
			}
		},
		5*time.Second,
	)
}
