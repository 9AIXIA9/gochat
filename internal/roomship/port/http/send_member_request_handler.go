package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/validator"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SendMemberRequestRequest struct {
	UserID   kernel.UserID   `json:"-" validate:"required"`
	RoomID   kernel.RoomID   `json:"room_id" validate:"required"`
	Password domain.Password `json:"password" validate:"max=100"`
	Content  string          `json:"content" validate:"max=100"`
}

func (r *SendMemberRequestRequest) Bind(ginContext *gin.Context) error {
	userID := ginutils.GetUserID(ginContext)
	r.UserID = userID
	return ginContext.ShouldBind(r)
}

func NewSendMemberRequestHandler(useCase application.SendMemberRequestUseCase, validator *validator.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *SendMemberRequestRequest) *application.SendMemberRequestInput {
			return &application.SendMemberRequestInput{
				UserID:   request.UserID,
				RoomID:   request.RoomID,
				Content:  request.Content,
				Password: request.Password,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.Response(ginContext, sharedHttp.ResponseSuccess)
		},
		func(ginContext *gin.Context, err error) {
			switch {
			default:
				zap.L().Error("SendMemberRequestHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
