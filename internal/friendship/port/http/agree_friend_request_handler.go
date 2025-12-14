package http

import (
	"errors"
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	ginutils "gochat/internal/infrastructure/gin"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AgreeFriendRequestRequest struct {
	UserID    kernel.UserID      `json:"-" validate:"required"`
	RequestID kernel.OperationID `uri:"request_id" validate:"required"`
}

func (r *AgreeFriendRequestRequest) Bind(ginContext *gin.Context) error {
	userID := ginutils.GetUserID(ginContext)
	r.UserID = userID
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

// NewAgreeFriendRequestHandler 同意好友请求
// @Summary      同意好友请求
// @Description  同意指定的好友请求，将对方加入好友列表
// @Tags         Friendship
// @Security     BearerAuth
// @Produce      json
// @Param        request_id  path      int                   true  "好友请求ID"
// @Success      200         {object}  sharedHttp.ApiResponse "同意成功"
// @Failure      400         {object}  sharedHttp.ApiResponse "请求参数错误或请求已处理"
// @Failure      401         {object}  sharedHttp.ApiResponse "未认证"
// @Failure      500         {object}  sharedHttp.ApiResponse "服务器内部错误"
// @Router       /friendship/request/{request_id}/agree [put]
func NewAgreeFriendRequestHandler(useCase application.AgreeFriendRequestUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *AgreeFriendRequestRequest) *application.AgreeFriendRequestInput {
			return &application.AgreeFriendRequestInput{
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
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "input is empty")
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "friend request not found")
			case errors.Is(err, domain.ErrFriendRequestNotForUser):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "friend request not for this user")
			case errors.Is(err, domain.ErrFriendRequestHasBeenHandled):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "friend request has been handled")
			default:
				zap.L().Error("AgreeFriendRequestHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.CodeServerError)
			}
		},
	)
}
