package http

import (
	"gochat/internal/friendship/application"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"

	_ "gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
)

type RefuseFriendRequestRequest struct {
	UserID    kernel.UserID      `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172190"`
	RequestID kernel.OperationID `uri:"request_id" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172191"`
}

func (r *RefuseFriendRequestRequest) Bind(ginContext *gin.Context) error {
	userID := ctxutil.UserIDFrom(ginContext.Request.Context())
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
// @Param        request_id  path   string   true  "好友请求ID"   example(019b593b-462e-74d6-bfda-0e103a172190)
// @Success      200         {object}  api.Response "拒绝成功"
// @Router       /friendship-requests/{request_id}/refuse [put]
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
	)
}
