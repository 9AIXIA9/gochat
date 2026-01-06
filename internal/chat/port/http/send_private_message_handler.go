package http

import (
	"gochat/internal/chat/application"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"

	_ "gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
)

type SendPrivateMessageRequest struct {
	SenderID    kernel.UserID `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172190"`
	RecipientID kernel.UserID `json:"recipient_id" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172190"`
	Content     string        `json:"content" validate:"required,max=1000" example:"Hello-Gochat!"`
}

func (r *SendPrivateMessageRequest) Bind(ginContext *gin.Context) error {
	senderID := ctxutil.UserIDFrom(ginContext.Request.Context())
	r.SenderID = senderID
	return ginContext.BindJSON(r)
}

// NewSendPrivateMessageHandler 发送私聊消息
// @Summary      发送私聊消息
// @Description  向指定用户发送一条私聊消息
// @Tags         Chat
// @Security     BearerAuth
// @Param        request  body      SendPrivateMessageRequest  true  "发送私聊消息请求体"
// @Success      200      {object}  api.Response     "发送成功"
// @Router       /chats/private-messages [post]
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
	)
}
