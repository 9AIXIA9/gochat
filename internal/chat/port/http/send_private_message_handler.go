package http

import (
	"errors"
	"gochat/internal/chat/application"
	ginutils "gochat/internal/infrastructure/gin"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SendPrivateMessageRequest struct {
	SenderID    kernel.UserID `json:"-" validate:"required"`
	RecipientID kernel.UserID `json:"recipient_id" validate:"required"`
	Content     string        `json:"content" validate:"required,max=1000"`
}

func (r *SendPrivateMessageRequest) Bind(ginContext *gin.Context) error {
	senderID := ginutils.GetUserID(ginContext)
	r.SenderID = senderID
	return ginContext.ShouldBind(r)
}

// NewSendPrivateMessageHandler 发送私聊消息
// @Summary      发送私聊消息
// @Description  向指定用户发送一条私聊消息
// @Tags         Chat
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      SendPrivateMessageRequest  true  "发送私聊消息请求体"
// @Success      201      {object}  sharedHttp.ApiResponse     "发送成功"
// @Failure      400      {object}  sharedHttp.ApiResponse     "请求参数错误"
// @Failure      401      {object}  sharedHttp.ApiResponse     "未认证"
// @Failure      500      {object}  sharedHttp.ApiResponse     "服务器内部错误"
// @Router       /chat/private [post]
func NewSendPrivateMessageHandler(useCase application.SendPrivateMessageUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *SendPrivateMessageRequest) *application.SendPrivateMessageInput {
			return &application.SendPrivateMessageInput{
				SenderID:    request.SenderID,
				RecipientID: request.RecipientID,
				Content:     request.Content,
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
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "recipientID not found")
			default:
				zap.L().Error("send private message handler failed", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.CodeServerError)
			}
		},
		5*time.Second,
	)
}
