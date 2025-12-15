package http

import (
	"errors"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	ginutils "gochat/internal/infrastructure/gin"
	sharedHttp "gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SendRoomMessageRequest struct {
	SenderID kernel.UserID `json:"-" validate:"required"`
	RoomID   kernel.RoomID `json:"room_id" validate:"required"`
	Content  string        `json:"content" validate:"required,max=1000"`
}

func (r *SendRoomMessageRequest) Bind(ginContext *gin.Context) error {
	senderID := ginutils.GetUserID(ginContext)
	r.SenderID = senderID
	return ginContext.ShouldBind(r)
}

// NewSendRoomMessageHandler 发送房间消息
// @Summary      发送房间消息
// @Description  向指定房间发送一条消息
// @Tags         Chat
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      SendRoomMessageRequest     true  "发送房间消息请求体"
// @Success      201      {object}  sharedHttp.Response     "发送成功"
// @Failure      400      {object}  sharedHttp.Response     "请求参数错误"
// @Failure      401      {object}  sharedHttp.Response     "未认证"
// @Failure      500      {object}  sharedHttp.Response     "服务器内部错误"
// @Router       /chat/room [post]
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
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "input is empty")
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "room not found")
			case errors.Is(err, domain.ErrNotMember):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "not belong to this room")
			case errors.Is(err, myErrors.ErrInvalidLength):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "invalid length")
			default:
				zap.L().Error("send room message handler failed", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.CodeServerError)
			}
		},
	)
}
