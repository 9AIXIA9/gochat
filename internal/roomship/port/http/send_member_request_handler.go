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

// NewSendMemberRequestHandler 发送入群请求
// @Summary      发送入群请求
// @Description  向指定房间发送入群请求，可携带验证信息及密码
// @Tags         Roomship
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      SendMemberRequestRequest  true  "发送入群请求体"
// @Success      201      {object}  api.Response    "发送成功"
// @Failure      400      {object}  api.Response    "请求参数错误或业务校验失败"
// @Failure      401      {object}  api.Response    "未认证"
// @Failure      500      {object}  api.Response    "服务器内部错误"
// @Router       /rooms/requests/ [post]
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
