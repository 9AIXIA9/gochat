package http

import (
	"errors"
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	ginutils "gochat/internal/infrastructure/gin"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SendFriendRequestRequest struct {
	FromID  kernel.UserID `json:"-" validate:"required"`
	ToID    kernel.UserID `json:"to_id" validate:"required"`
	Content string        `json:"content" validate:"required,max=100"`
}

func (r *SendFriendRequestRequest) Bind(ginContext *gin.Context) error {
	fromID := ginutils.GetUserID(ginContext)
	r.FromID = fromID
	return ginContext.ShouldBind(r)
}

func NewSendFriendRequestHandler(useCase application.SendFriendRequestUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *SendFriendRequestRequest) *application.SendFriendRequestInput {
			return &application.SendFriendRequestInput{
				FromID:  request.FromID,
				ToID:    request.ToID,
				Content: request.Content,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.ResponseSuccess(ginContext)
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "input is empty")
			case errors.Is(err, domain.ErrAddYourselfAsFriend):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "you cannot add yourself as a friend")
			case errors.Is(err, domain.ErrAlreadyBeenFriends):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "you are already friends")
			case errors.Is(err, domain.ErrFriendRequestExists):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "friend request already sent")
			case errors.Is(err, domain.ErrFriendRequestContentTooLong):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "friend request content too long")
			default:
				zap.L().Error("SendFriendRequestHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.CodeServerError)
			}
		},
		5*time.Second,
	)
}
