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
	OperatorID kernel.UserID    `json:"-" validate:"required" example:"019b5929-5f6a-73aa-8adf-57cebe980725"`
	UserID     kernel.UserID    `uri:"user_id" validate:"required" example:"019b5929-65bc-7549-89c6-3f7dc6872577"`
	BaseID     kernel.MessageID `form:"base_id"  validate:"omitempty" example:"019b593b-462e-74d6-bfda-0e103a172190"` // 用于分页游标
	Limit      int              `form:"limit" validate:"omitempty,min=1,max=100" example:"50"`                        // 每页条数
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
// @Param        user_id  path      string  true  "私聊对象用户ID"  example(019b5929-65bc-7549-89c6-3f7dc6872577)
// @Param        base_id  query     string  false "分页游标"       example(019b593b-462e-74d6-bfda-0e103a172190)
// @Param        limit    query     int     false "分页大小，默认 20，最大 100"      minimum(1) maximum(100) default(20) example(50)
// @Success      200      {object}  api.Response{data=ListPrivateMessagesResponseData} "成功返回私聊消息记录"
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
