package http

import (
	"errors"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	ginutils "gochat/internal/infrastructure/gin"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

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
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeNotBelongTo, "not belong to this room")
			case errors.Is(err, myErrors.ErrInvalidLength):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidRoom, "invalid room")
			default:
				zap.L().Error("send room message handler failed", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.CodeServerError)
			}
		},
		5*time.Second,
	)
}
