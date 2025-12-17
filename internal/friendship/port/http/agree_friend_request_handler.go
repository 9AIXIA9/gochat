package http

import (
	"gochat/internal/friendship/application"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

type AgreeFriendRequestRequest struct {
	UserID    kernel.UserID      `json:"-" validate:"required"`
	RequestID kernel.OperationID `uri:"request_id" validate:"required"`
}

func (r *AgreeFriendRequestRequest) Bind(ginContext *gin.Context) error {
	userID := ginutils.GetUserID(ginContext)
	r.UserID = userID
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

// NewAgreeFriendRequestHandler 同意好友请求
// @Summary      同意好友请求
// @Description  同意指定的好友请求，将对方加入好友列表
// @Tags         Friendship
// @Security     BearerAuth
// @Produce      json
// @Param        request_id  path      kernel.OperationID                   true  "好友请求ID"
// @Success      200         {object}  api.Response "同意成功"
// @Failure      400         {object}  api.Response "请求参数错误或请求已处理"
// @Failure      401         {object}  api.Response "未认证"
// @Failure      500         {object}  api.Response "服务器内部错误"
// @Router       /friendship/request/{request_id}/agree [put]
func NewAgreeFriendRequestHandler(useCase application.AgreeFriendRequestUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *AgreeFriendRequestRequest) *application.AgreeFriendRequestInput {
			return &application.AgreeFriendRequestInput{
				UserID:    request.UserID,
				RequestID: request.RequestID,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.ResponseSuccess(ginContext)
		},
	)
}
