package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/validator"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
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
			case errors.Is(err, domain.ErrMemberRequestAlreadyExists):
				sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "you have already sent a member request to this room")
			case errors.Is(err, domain.ErrIsAlreadyMember):
				sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "you are already a member of this room")
			case errors.Is(err, myErrors.ErrNotFound):
				sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "the room does not exist")
			case errors.Is(err, domain.ErrInvalidPassword):
				sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "the password is incorrect")
			default:
				zap.L().Error("SendMemberRequestHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
