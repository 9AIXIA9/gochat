package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	"gochat/internal/shared/kernel"

	_ "gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
)

type AgreeMemberRequestRequest struct {
	UserID    kernel.UserID      `json:"-" validate:"required"`
	RequestID kernel.OperationID `uri:"request_id" validate:"required"`
}

func (r *AgreeMemberRequestRequest) Bind(ginContext *gin.Context) error {
	userID := ginutils.GetUserID(ginContext)
	r.UserID = userID
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

// NewAgreeMemberRequestHandler 同意入群请求
// @Summary      同意入群请求
// @Description  同意指定房间成员请求，将对方加入房间
// @Tags         Roomship
// @Security     BearerAuth
// @Produce      json
// @Param        request_id  path      kernel.OperationID                   true  "成员请求ID"
// @Success      200         {object}  api.Response "同意成功"
// @Failure      400         {object}  api.Response "请求参数错误或请求已处理/无权限"
// @Failure      401         {object}  api.Response "未认证"
// @Failure      500         {object}  api.Response "服务器内部错误"
// @Router       /rooms/requests/{request_id}/agree [put]
func NewAgreeMemberRequestHandler(useCase application.AgreeMemberRequestUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *AgreeMemberRequestRequest) *application.AgreeMemberRequestInput {
			return &application.AgreeMemberRequestInput{
				UserID:    request.UserID,
				RequestID: request.RequestID,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.ResponseSuccess(ginContext)
		},
	)
}
