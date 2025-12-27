package http

import (
	"gochat/internal/friendship/application"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/kernel"

	_ "gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
)

type SendFriendRequestRequest struct {
	FromID  kernel.UserID `json:"-" validate:"required"`
	ToID    kernel.UserID `json:"to_id" validate:"required"`
	Content string        `json:"content" validate:"required,max=100"`
}

func (r *SendFriendRequestRequest) Bind(ginContext *gin.Context) error {
	fromID := ginutils.GetUserID(ginContext)
	r.FromID = fromID
	return ginContext.BindJSON(r)
}

// NewSendFriendRequestHandler 发送好友请求
// @Summary      发送好友请求
// @Description  向指定用户发送好友请求
// @Tags         Friendship
// @Security     BearerAuth
// @Param        request  body      SendFriendRequestRequest  true  "发送好友请求体"
// @Success      201      {object}  api.Response    "发送成功"
// @Failure      400      {object}  api.Response    "请求参数错误或业务校验失败"
// @Failure      401      {object}  api.Response    "未认证"
// @Failure      500      {object}  api.Response    "服务器内部错误"
// @Router       /friendship-requests/ [post]
func NewSendFriendRequestHandler(useCase application.SendFriendRequestUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *SendFriendRequestRequest) *application.SendFriendRequestInput {
			return &application.SendFriendRequestInput{
				FromID:  request.FromID,
				ToID:    request.ToID,
				Content: request.Content,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.ResponseSuccess(ginContext)
		},
	)
}
