package http

import (
	"errors"
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RefuseFriendRequestRequest struct {
	UserID    kernel.UserID      `json:"-" validate:"required"`
	RequestID kernel.OperationID `uri:"request_id" validate:"required"`
}

func (r *RefuseFriendRequestRequest) Bind(ginContext *gin.Context) error {
	userID := ginutils.GetUserID(ginContext)
	r.UserID = userID
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

// NewRefuseFriendRequestHandler 拒绝好友请求
// @Summary      拒绝好友请求
// @Description  拒绝指定的好友请求
// @Tags         Friendship
// @Security     BearerAuth
// @Produce      json
// @Param        request_id  path      int                   true  "好友请求ID"
// @Success      200         {object}  api.Response "拒绝成功"
// @Failure      400         {object}  api.Response "请求参数错误或请求已处理"
// @Failure      401         {object}  api.Response "未认证"
// @Failure      500         {object}  api.Response "服务器内部错误"
// @Router       /friendship/request/{request_id}/refuse [put]
func NewRefuseFriendRequestHandler(useCase application.RefuseFriendRequestUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *RefuseFriendRequestRequest) *application.RefuseFriendRequestInput {
			return &application.RefuseFriendRequestInput{
				UserID:    request.UserID,
				RequestID: request.RequestID,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.ResponseSuccess(ginContext)
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.ResponseWithMessage(ginContext, api.CodeInvalidParam, "input is empty")
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.ResponseWithMessage(ginContext, api.CodeInvalidParam, "friend request not found")
			case errors.Is(err, domain.ErrFriendRequestNotForUser):
				ginutils.ResponseWithMessage(ginContext, api.CodeInvalidParam, "friend request not for this user")
			case errors.Is(err, domain.ErrFriendRequestHasBeenHandled):
				ginutils.ResponseWithMessage(ginContext, api.CodeInvalidParam, "friend request has been handled")
			default:
				zap.L().Error("RefuseFriendRequestHandler error", zap.Error(err))
				ginutils.Response(ginContext, api.CodeServerError)
			}
		},
	)
}
