package http

import (
	"gochat/internal/friendship/application"
	ginutils "gochat/internal/infrastructure/gin"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

type AgreeFriendRequestRequest struct {
	UserID    kernel.UserID      `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172190"`
	RequestID kernel.OperationID `uri:"request_id" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172191"`
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
// @Param        request_id  path   string   true  "好友请求ID"   example(019b593b-462e-74d6-bfda-0e103a172190)
// @Success      200         {object}  api.Response "同意成功"
// @Router       /friendship-requests/{request_id}/agree [put]
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
	)
}
