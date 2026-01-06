package http

import (
	"gochat/internal/chat/application"
	ginutils "gochat/internal/infrastructure/gin"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"

	"github.com/gin-gonic/gin"
)

type ReadPrivateMessagesRequest struct {
	SenderID    kernel.UserID `json:"sender_id" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172190"`
	RecipientID kernel.UserID `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172191"`
}

func (r *ReadPrivateMessagesRequest) Bind(ginContext *gin.Context) error {
	r.RecipientID = ctxutil.UserIDFrom(ginContext.Request.Context())
	return ginContext.BindJSON(r)
}

// NewReadPrivateMessagesHandler 阅读指定用户消息
// @Summary      阅读指定用户消息
// @Description  阅读来自指定用户的所有未读消息
// @Tags         Chat
// @Security     BearerAuth
// @Param        request  body      ReadPrivateMessagesRequest     true  "阅读指定用户消息请求体"
// @Success      200      {object}  api.Response     "成功已读"
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
