package http

import (
	"gochat/internal/chat/application"
	ginutils "gochat/internal/infrastructure/gin"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

type SendRoomMessageRequest struct {
	SenderID kernel.UserID `json:"-" validate:"required"`
	RoomID   kernel.RoomID `json:"room_id" validate:"required"`
	Content  string        `json:"content" validate:"required,max=1000"`
}

func (r *SendRoomMessageRequest) Bind(ginContext *gin.Context) error {
	senderID := ginutils.GetUserID(ginContext)
	r.SenderID = senderID
	return ginContext.BindJSON(r)
}

// NewSendRoomMessageHandler 发送房间消息
// @Summary      发送房间消息
// @Description  向指定房间发送一条消息
// @Tags         Chat
// @Security     BearerAuth
// @Param        request  body      SendRoomMessageRequest     true  "发送房间消息请求体"
// @Success      200      {object}  api.Response     "发送成功"
// @Router       /chats/rooms/messages [post]
func NewSendRoomMessageHandler(useCase application.SendRoomMessageUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *SendRoomMessageRequest) *application.SendRoomMessageInput {
			return &application.SendRoomMessageInput{
				SenderID: request.SenderID,
				RoomID:   request.RoomID,
				Content:  request.Content,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.ResponseSuccess(ginContext)
		},
	)
}
