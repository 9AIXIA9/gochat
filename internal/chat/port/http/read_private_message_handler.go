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
	return ginContext.BindJSON(r)
}

// NewReadPrivateMessagesHandler 发送房间消息
// @Summary      发送房间消息
// @Description  向指定房间发送一条消息
// @Tags         Chat
// @Security     BearerAuth
// @Param        request  body      ReadPrivateMessagesRequest     true  "发送房间消息请求体"
// @Success      200      {object}  api.Response     "发送成功"
// @Router       /chats/private-messages/read [put]
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
