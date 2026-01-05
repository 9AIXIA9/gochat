package http

import (
	"gochat/internal/chat/application"
	ginutils "gochat/internal/infrastructure/gin"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"

	"github.com/gin-gonic/gin"
)

type ReadRoomMessagesRequest struct {
	RoomID kernel.RoomID `json:"room_id" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172190"`
	UserID kernel.UserID `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172191"`
}

func (r *ReadRoomMessagesRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ctxutil.UserIDFrom(ginContext.Request.Context())
	return ginContext.BindJSON(r)
}

// NewReadRoomMessagesHandler 阅读指定房间消息
// @Summary      阅读指定房间消息
// @Description  阅读指定房间的所有消息
// @Tags         Chat
// @Security     BearerAuth
// @Param        request  body      ReadRoomMessagesRequest     true  "阅读指定房间消息请求体"
// @Success      200      {object}  api.Response     "成功已读"
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
