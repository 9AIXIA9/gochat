package http

import (
	"gochat/internal/friendship/application"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/kernel"

	_ "gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
)

type SendFriendRequestRequest struct {
	FromID  kernel.UserID `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172190"`
	ToID    kernel.UserID `json:"to_id" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	Content string        `json:"content" validate:"required,max=100" example:"019b593b-462e-74d6-bfda-0e103a172192"`
}

func (r *SendFriendRequestRequest) Bind(ginContext *gin.Context) error {
	fromID := ginutils.GetUserID(ginContext)
	r.FromID = fromID
	return ginContext.BindJSON(r)
}

// NewSendFriendRequestHandler 发送好友请求
// @Summary      发送好友请求
// @Description  向指定用户发送好友请求
// @Tags         Friendship
// @Security     BearerAuth
// @Param        request  body      SendFriendRequestRequest  true  "发送好友请求体"
// @Success      200      {object}  api.Response    "发送成功"
// @Router       /friendship-requests [post]
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
	)
}
