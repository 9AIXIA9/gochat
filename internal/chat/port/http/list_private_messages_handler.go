package http

import (
	"errors"
	"gochat/internal/chat/application"
	"gochat/internal/chat/dto"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const defaultPrivateMessagesLimit = 20

type ListPrivateMessagesRequest struct {
	UserID kernel.UserID    `json:"-" validate:"required"`
	BaseID kernel.MessageID `form:"base_id"`                                  // 用于分页游标
	Limit  int              `form:"limit" validate:"omitempty,min=1,max=100"` // 每页条数
}

func (r *ListPrivateMessagesRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	// 绑定查询参数
	if err := ginContext.ShouldBindQuery(r); err != nil {
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

// NewListPrivateMessagesHandler 获取私聊消息列表
// @Summary      获取私聊消息列表
// @Description  获取当前登录用户的私聊消息列表，可基于 base_id 游标和 limit 分页
// @Tags         Chat
// @Security     BearerAuth
// @Produce      json
// @Param        base_id  query     kernel.MessageID    false "分页游标，返回该ID之前的消息"
// @Param        limit    query     int    false "分页大小，默认20，最大100"
// @Success      200      {object}  api.Response{data=ListPrivateMessagesResponseData} "成功返回私聊消息列表"
// @Failure      400      {object}  api.Response "请求参数错误"
// @Failure      401      {object}  api.Response "未认证"
// @Failure      500      {object}  api.Response "服务器内部错误"
// @Router       /chat/private [get]
func NewListPrivateMessagesHandler(useCase application.ListPrivateMessagesUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *ListPrivateMessagesRequest) *application.ListPrivateMessagesInput {
			return &application.ListPrivateMessagesInput{
				UserID: request.UserID,
				BaseID: request.BaseID,
				Limit:  request.Limit,
			}
		},
		func(ginContext *gin.Context, output *application.ListPrivateMessagesOutput) {
			ginutils.ResponseSuccessWithData(ginContext, &ListPrivateMessagesResponseData{
				PrivateMessages: dto.ToPrivateMessageDTOs(output.PrivateMessages),
			})
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.ResponseWithMessage(ginContext, api.CodeInvalidParam, "input is empty")
			default:
				zap.L().Error("ListPrivateMessageHandler error", zap.Error(err))
				ginutils.Response(ginContext, api.CodeServerError)
			}
		},
	)
}
