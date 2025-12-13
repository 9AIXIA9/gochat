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

type UpdateUserProfileRequest struct {
	UserID      kernel.UserID      `json:"-" validate:"required"`
	Name        string             `json:"name"`
	Gender      kernel.Gender      `json:"gender"`
	Email       kernel.Email       `json:"email"`
	PhoneNumber kernel.PhoneNumber `json:"phone_number"`
	Address     kernel.Address     `json:"address"`
	Sign        string             `json:"sign"`
}

func (r *UpdateUserProfileRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	// 绑定查询参数
	if err := ginContext.ShouldBind(r); err != nil {
		return err
	}
	return nil
}

// NewUpdateUserProfileHandler 更新当前用户资料
// @Summary      更新当前用户资料
// @Description  更新当前登录用户的基础资料（昵称、邮箱、电话等）
// @Tags         Profile
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      UpdateUserProfileRequest  true  "更新用户资料请求体"
// @Success      200      {object}  sharedHttp.ApiResponse    "更新成功"
// @Failure      400      {object}  sharedHttp.ApiResponse    "请求参数错误"
// @Failure      401      {object}  sharedHttp.ApiResponse    "未认证"
// @Failure      404      {object}  sharedHttp.ApiResponse    "用户资料不存在"
// @Failure      500      {object}  sharedHttp.ApiResponse    "服务器内部错误"
// @Router       /profile/me [put]
func NewUpdateUserProfileHandler(useCase application.UpdateUserProfileUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *UpdateUserProfileRequest) *application.UpdateUserProfileInput {
			return &application.UpdateUserProfileInput{
				UserID:      request.UserID,
				Name:        request.Name,
				Gender:      request.Gender,
				Email:       request.Email,
				PhoneNumber: request.PhoneNumber,
				Address:     request.Address,
				Sign:        request.Sign,
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
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeNotFound, "user profile not found")
			default:
				zap.L().Error("UpdateUserProfileHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.CodeServerError)
			}
		},
		5*time.Second,
	)
}
