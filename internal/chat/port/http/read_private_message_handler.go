package http

import (
	"gochat/internal/chat/application"
	ginutils "gochat/internal/infrastructure/gin"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

type ReadPrivateMessagesRequest struct {
	SenderID    kernel.UserID `json:"sender_id" validate:"required"`
	RecipientID kernel.UserID `json:"-" validate:"required"`
}

func (r *ReadPrivateMessagesRequest) Bind(ginContext *gin.Context) error {
	r.RecipientID = ginutils.GetUserID(ginContext)
	return ginContext.ShouldBind(r)
}

// NewReadPrivateMessagesHandler 发送房间消息
// @Summary      发送房间消息
// @Description  向指定房间发送一条消息
// @Tags         Chat
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      ReadPrivateMessagesRequest     true  "发送房间消息请求体"
// @Success      201      {object}  api.Response     "发送成功"
// @Failure      400      {object}  api.Response     "请求参数错误"
// @Failure      401      {object}  api.Response     "未认证"
// @Failure      500      {object}  api.Response     "服务器内部错误"
// @Router       /chat/room [post]
func NewReadPrivateMessagesHandler(useCase application.ReadPrivateMessagesUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *ReadPrivateMessagesRequest) *application.ReadPrivateMessagesInput {
			return &application.ReadPrivateMessagesInput{
				SenderID:    request.SenderID,
				RecipientID: request.RecipientID,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.ResponseSuccess(ginContext)
		},
	)
}
