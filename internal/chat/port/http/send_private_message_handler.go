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
