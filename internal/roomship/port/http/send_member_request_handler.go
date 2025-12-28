package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/kernel"

	_ "gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
)

type SendMemberRequestRequest struct {
	UserID   kernel.UserID   `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	RoomID   kernel.RoomID   `json:"room_id" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172192"`
	Password domain.Password `json:"password" validate:"max=100" example:"secret"`
	Content  string          `json:"content" validate:"max=100" example:"I would like to join the room."`
}

func (r *SendMemberRequestRequest) Bind(ginContext *gin.Context) error {
	userID := ginutils.GetUserID(ginContext)
	r.UserID = userID
	return ginContext.BindJSON(r)
}

// NewSendMemberRequestHandler 发送入群请求
// @Summary      发送入群请求
// @Description  向指定房间发送入群请求，可携带验证信息及密码
// @Tags         Roomship
// @Security     BearerAuth
// @Param        request  body      SendMemberRequestRequest  true  "发送入群请求体"
// @Success      200      {object}  api.Response    "发送成功"
// @Router       /rooms/requests [post]
func NewSendMemberRequestHandler(useCase application.SendMemberRequestUseCase, validator ginutils.Validator) gin.HandlerFunc {
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
			ginutils.ResponseSuccess(ginContext)
		},
	)
}
