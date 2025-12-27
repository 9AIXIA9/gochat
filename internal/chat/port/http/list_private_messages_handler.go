package http

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/dto"
	ginutils "gochat/internal/infrastructure/gin"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

const defaultPrivateMessagesLimit = 20

type ListPrivateMessagesRequest struct {
	OperatorID kernel.UserID    `json:"-" validate:"required"`
	UserID     kernel.UserID    `uri:"user_id" validate:"required"`
	BaseID     kernel.MessageID `form:"base_id"`                                  // 用于分页游标
	Limit      int              `form:"limit" validate:"omitempty,min=1,max=100"` // 每页条数
}

func (r *ListPrivateMessagesRequest) Bind(ginContext *gin.Context) error {
	r.OperatorID = ginutils.GetUserID(ginContext)
	// 绑定查询参数
	if err := ginContext.ShouldBindQuery(r); err != nil {
		return err
	}
	if err := ginContext.BindUri(r); err != nil {
		return err
	}
	// 默认值
	if r.Limit == 0 {
		r.Limit = defaultPrivateMessagesLimit
	}
	return nil
}

type ListPrivateMessagesResponseData struct {
	PrivateMessages []*dto.PrivateMessage `json:"private_messages,omitempty"`
}

// NewListPrivateMessagesHandler 获取用户聊天记录
// @Summary      获取用户聊天记录
// @Description  获取用户和对方聊天记录，可基于 base_id 游标和 limit 分页
// @Tags         Chat
// @Security     BearerAuth
// @Param        user_id  path      kernel.UserID       true  "私聊对象用户ID"
// @Param        base_id  query     kernel.MessageID    false "分页游标，返回该ID之前的消息"
// @Param        limit    query     int    false "分页大小，默认20，最大100"
// @Success      200      {object}  api.Response{data=ListPrivateMessagesResponseData} "成功返回私聊消息记录"
// @Failure      400      {object}  api.Response "请求参数错误"
// @Failure      401      {object}  api.Response "未认证"
// @Failure      500      {object}  api.Response "服务器内部错误"
// @Router       /chats/private-messages/{user_id} [get]
func NewListPrivateMessagesHandler(useCase application.ListPrivateMessagesUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *ListPrivateMessagesRequest) *application.ListPrivateMessagesInput {
			return &application.ListPrivateMessagesInput{
				OperatorID: request.OperatorID,
				UserID:     request.UserID,
				BaseID:     request.BaseID,
				Limit:      request.Limit,
			}
		},
		func(ginContext *gin.Context, output *application.ListPrivateMessagesOutput) {
			ginutils.ResponseSuccessWithData(ginContext, &ListPrivateMessagesResponseData{
				PrivateMessages: dto.ToPrivateMessageDTOs(output.PrivateMessages),
			})
		},
	)
}
