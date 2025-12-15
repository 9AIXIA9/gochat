package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/profile/application"
	sharedHttp "gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

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

// NewUpdateRoomProfileHandler 更新房间资料
// @Summary      更新房间资料
// @Description  更新指定房间的资料（名称、简介等），通常需要房主或管理员权限
// @Tags         Profile
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      UpdateRoomProfileRequest  true  "更新房间资料请求体"
// @Success      200      {object}  sharedHttp.Response    "更新成功"
// @Failure      400      {object}  sharedHttp.Response    "请求参数错误"
// @Failure      401      {object}  sharedHttp.Response    "未认证"
// @Failure      404      {object}  sharedHttp.Response    "房间资料不存在"
// @Failure      500      {object}  sharedHttp.Response    "服务器内部错误"
// @Router       /profile/room [put]
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
	)
}
