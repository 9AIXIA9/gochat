package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/notification/application"
	"gochat/internal/notification/dto"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

const defaultSystemMessagesLimit = 20

type ListSystemMessagesRequest struct {
	UserID kernel.UserID    `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	BaseID kernel.MessageID `form:"base_id"  validate:"omitempty" example:"019b593b-462e-74d6-bfda-0e103a172192"` // 用于分页游标
	Limit  int              `form:"limit" validate:"omitempty,min=1,max=100" example:"20"`                        // 每页条数
}

func (r *ListSystemMessagesRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	// 绑定查询参数
	if err := ginContext.ShouldBindQuery(r); err != nil {
		return err
	}
	// 默认值
	if r.Limit == 0 {
		r.Limit = defaultSystemMessagesLimit
	}
	return nil
}

type ListSystemMessagesResponseData struct {
	SystemMessages []*dto.SystemMessage `json:"system_messages,omitempty"`
}

// NewListSystemMessagesHandler 获取系统通知列表
// @Summary      获取系统通知列表
// @Description  获取当前登录用户的系统通知列表，可基于 base_id 游标和 limit 分页
// @Tags         Notification
// @Security     BearerAuth
// @Param        base_id  query     string    false "分页游标，返回该ID之前的消息"    example("019b593b-462e-74d6-bfda-0e103a172190")
// @Param        limit    query     int    false "分页大小，默认20，最大100"   minimum(1) maximum(100) default(20) example(50)
// @Success      200      {object}  api.Response{data=ListSystemMessagesResponseData} "成功返回系统通知列表"
// @Router       /notifications/system-messages [get]
func NewListSystemMessagesHandler(useCase application.ListSystemMessagesUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *ListSystemMessagesRequest) *application.ListSystemMessagesInput {
			return &application.ListSystemMessagesInput{
				UserID: request.UserID,
				BaseID: request.BaseID,
				Limit:  request.Limit,
			}
		},
		func(ginContext *gin.Context, output *application.ListSystemMessagesOutput) {
			ginutils.ResponseSuccessWithData(ginContext, &ListSystemMessagesResponseData{
				SystemMessages: dto.ToSystemMessageDTOs(output.SystemMessages),
			})
		},
	)
}
