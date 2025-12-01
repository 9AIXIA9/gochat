package http

import (
	"gochat/internal/friendship/application"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/validator"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
)

type SendFriendRequestRequest struct {
	FromID   kernel.UserID     `json:"-" validate:"required"`
	ToNumber kernel.UserNumber `json:"to_number" validate:"required"`
	Content  string            `json:"content" validate:"required,max=100"`
}

func (r *SendFriendRequestRequest) Bind(ginContext *gin.Context) error {
	fromID := ginutils.GetUserID(ginContext)
	r.FromID = fromID
	return ginContext.ShouldBind(r)
}

func NewSendFriendRequestHandler(useCase application.SendFriendRequestUseCase, validator *validator.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *SendFriendRequestRequest) *application.SendFriendRequestInput {
			return &application.SendFriendRequestInput{
				FromID:   request.FromID,
				ToNumber: request.ToNumber,
				Content:  request.Content,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.Response(ginContext, sharedHttp.ResponseSuccess)
		},
		func(ginContext *gin.Context, err error) {
			switch err.(type) {
			//TODO handle specific errors
			default:
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
