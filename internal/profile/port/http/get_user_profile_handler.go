package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/profile/application"
	"gochat/internal/profile/dto"
	sharedHttp "gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetUserProfileRequest struct {
	UserID kernel.UserID `uri:"user_id" validate:"required"`
}

func (r *GetUserProfileRequest) Bind(ginContext *gin.Context) error {
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

type GetUserProfileResponseData struct {
	Profile *dto.UserProfile `json:"profile"`
}

// NewGetUserProfileHandler 获取用户资料
// @Summary      获取用户资料
// @Description  根据用户ID获取用户公开资料
// @Tags         Profile
// @Produce      json
// @Param        user_id  path      int                         true  "用户ID"
// @Success      200      {object}  GetUserProfileResponseData  "成功返回用户资料"
// @Failure      400      {object}  sharedHttp.Response "请求参数错误"
// @Failure      404      {object}  sharedHttp.Response "用户不存在"
// @Failure      500      {object}  sharedHttp.Response "服务器内部错误"
// @Router       /profile/user/{user_id} [get]
func NewGetUserProfileHandler(useCase application.GetUserProfileUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *GetUserProfileRequest) *application.GetUserProfileInput {
			return &application.GetUserProfileInput{
				UserID: request.UserID,
			}
		},
		func(ginContext *gin.Context, output *application.GetUserProfileOutput) {
			ginutils.ResponseSuccessWithData(ginContext, &GetUserProfileResponseData{
				Profile: dto.ToUserProfileDTO(output.Profile),
			})
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "input is empty")
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeNotFound, "user is not found")
			default:
				zap.L().Error("GetUserProfileHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.CodeServerError)
			}
		},
	)
}
