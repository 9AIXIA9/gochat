package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RefuseMemberRequestRequest struct {
	UserID    kernel.UserID      `json:"-" validate:"required"`
	RequestID kernel.OperationID `uri:"request_id" validate:"required"`
}

func (r *RefuseMemberRequestRequest) Bind(ginContext *gin.Context) error {
	userID := ginutils.GetUserID(ginContext)
	r.UserID = userID
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

// NewRefuseMemberRequestHandler 拒绝入群请求
// @Summary      拒绝入群请求
// @Description  拒绝指定房间成员请求
// @Tags         Roomship
// @Security     BearerAuth
// @Produce      json
// @Param        request_id  path      int                   true  "成员请求ID"
// @Success      200         {object}  sharedHttp.ApiResponse "拒绝成功"
// @Failure      400         {object}  sharedHttp.ApiResponse "请求参数错误或请求已处理/无权限"
// @Failure      401         {object}  sharedHttp.ApiResponse "未认证"
// @Failure      500         {object}  sharedHttp.ApiResponse "服务器内部错误"
// @Router       /roomship/request/{request_id}/refuse [put]
func NewRefuseMemberRequestHandler(useCase application.RefuseMemberRequestUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *RefuseMemberRequestRequest) *application.RefuseMemberRequestInput {
			return &application.RefuseMemberRequestInput{
				UserID:    request.UserID,
				RequestID: request.RequestID,
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
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "request is not found")
			case errors.Is(err, domain.ErrNotAdmin):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "you have no permission to refuse this request")
			default:
				zap.L().Error("RefuseMemberRequestHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.CodeServerError)
			}
		},
	)
}
