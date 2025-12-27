package http

import (
	"gochat/internal/chat/application"
	ginutils "gochat/internal/infrastructure/gin"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

type ReadRoomMessagesRequest struct {
	RoomID kernel.RoomID `json:"room_id" validate:"required"`
	UserID kernel.UserID `json:"-" validate:"required"`
}

func (r *ReadRoomMessagesRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	return ginContext.BindJSON(r)
}

// NewReadRoomMessagesHandler 发送房间消息
// @Summary      发送房间消息
// @Description  向指定房间发送一条消息
// @Tags         Chat
// @Security     BearerAuth
// @Param        request  body      ReadRoomMessagesRequest     true  "发送房间消息请求体"
// @Success      200      {object}  api.Response     "发送成功"
// @Router       /chats/rooms/messages/read [put]
func NewReadRoomMessagesHandler(useCase application.ReadRoomMessagesUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *ReadRoomMessagesRequest) *application.ReadRoomMessagesInput {
			return &application.ReadRoomMessagesInput{
				UserID: request.UserID,
				RoomID: request.RoomID,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.ResponseSuccess(ginContext)
		},
	)
}
