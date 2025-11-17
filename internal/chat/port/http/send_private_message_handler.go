package http

import (
	"errors"
	"gochat/internal/chat/application/usecase"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/validator"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SendPrivateMessageRequest struct {
	SenderID        kernel.UserID     `json:"-" validate:"required"`
	RecipientNumber kernel.UserNumber `json:"recipient_number" validate:"required"`
	Content         string            `json:"content" validate:"required,max=1000"`
}

func (r *SendPrivateMessageRequest) Bind(ginContext *gin.Context) error {
	senderID := ginutils.GetUserID(ginContext)
	r.SenderID = senderID
	return ginContext.ShouldBind(r)
}

func NewSendPrivateMessageHandler(useCase usecase.SendPrivateMessageUseCase, validator *validator.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *SendPrivateMessageRequest) *usecase.SendPrivateMessageInput {
			return &usecase.SendPrivateMessageInput{
				SenderID:        request.SenderID,
				RecipientNumber: request.RecipientNumber,
				Content:         request.Content,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.Response(ginContext, sharedHttp.ResponseSuccess)
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "input is empty"))
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "recipient not found"))
			default:
				zap.L().Error("send private message handler failed", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
