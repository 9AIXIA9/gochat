package http

import (
	"errors"
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	ginutils "gochat/internal/infrastructure/gin"
	sharedHttp "gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

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

// NewSendFriendRequestHandler 发送好友请求
// @Summary      发送好友请求
// @Description  向指定用户发送好友请求
// @Tags         Friendship
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      SendFriendRequestRequest  true  "发送好友请求体"
// @Success      201      {object}  sharedHttp.Response    "发送成功"
// @Failure      400      {object}  sharedHttp.Response    "请求参数错误或业务校验失败"
// @Failure      401      {object}  sharedHttp.Response    "未认证"
// @Failure      500      {object}  sharedHttp.Response    "服务器内部错误"
// @Router       /friendship/request [post]
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
	)
}
