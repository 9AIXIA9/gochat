package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
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
// @Param        request_id  path      kernel.OperationID                   true  "成员请求ID"
// @Success      200         {object}  api.Response "拒绝成功"
// @Failure      400         {object}  api.Response "请求参数错误或请求已处理/无权限"
// @Failure      401         {object}  api.Response "未认证"
// @Failure      500         {object}  api.Response "服务器内部错误"
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
	)
}
